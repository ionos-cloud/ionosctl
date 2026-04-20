package object

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
)

const flagBypassGovernanceRetention = "bypass-governance-retention"

func DeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete an object or all objects from a bucket",
		LongDesc:  "Delete a single object by key, or all objects (including versions and delete markers) from a bucket using --all.",
		Example:   "ionosctl object-storage object delete --name my-bucket --key photos/image.jpg\nionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id abc123 -f\nionosctl object-storage object delete --name my-bucket --all -f\nionosctl object-storage object delete --name my-bucket --all --bypass-governance-retention -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagName, flagKey},
				[]string{constants.FlagName, constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			bypassGovernance := viper.GetBool(core.GetFlagName(c.NS, flagBypassGovernanceRetention))

			if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
				if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete ALL objects in bucket %q", name), viper.GetBool(constants.ArgForce)) {
					return fmt.Errorf(confirm.UserDenied)
				}

				if err := emptyBucket(c, client.MustObjectStorage().ObjectStorageClient, name, bypassGovernance); err != nil {
					return err
				}

				c.Msg("All objects deleted from bucket %q", name)
				return nil
			}

			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete object %q from bucket %q", key, name), viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			req := client.MustObjectStorage().ObjectStorageClient.ObjectsApi.DeleteObject(c.Context, name, key).XAmzBypassGovernanceRetention(bypassGovernance)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			_, _, err := req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Object %q deleted from bucket %q\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key to delete",
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID to delete a specific version")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all objects in the bucket")
	cmd.AddBoolFlag(flagBypassGovernanceRetention, "", false, "Bypass Governance-mode Object Lock restrictions to delete the object")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

// emptyBucket deletes all objects, object versions, and delete markers in the bucket.
func emptyBucket(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string, bypassGovernance bool) error {
	// First pass: delete current objects via ListObjectsV2.
	if err := deleteCurrentObjects(c, s3, bucket, bypassGovernance); err != nil {
		return err
	}

	// Second pass: delete all versions and delete markers via ListObjectVersions.
	if err := deleteAllVersions(c, s3, bucket, bypassGovernance); err != nil {
		return err
	}

	return nil
}

func deleteCurrentObjects(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string, bypassGovernance bool) error {
	var continuationToken string
	totalDeleted := 0

	for {
		req := s3.ObjectsApi.ListObjectsV2(c.Context, bucket).MaxKeys(1000)
		if continuationToken != "" {
			req = req.ContinuationToken(continuationToken)
		}

		result, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("listing objects: %w", err)
		}

		if len(result.Contents) == 0 {
			break
		}

		if err := batchDelete(c.Context, s3, bucket, objectsToIdentifiers(result.Contents), bypassGovernance); err != nil {
			return err
		}

		totalDeleted += len(result.Contents)
		c.Verbose("Deleted %d objects...", totalDeleted)

		if !result.IsTruncated {
			break
		}

		if result.NextContinuationToken != nil {
			continuationToken = *result.NextContinuationToken
		} else {
			break
		}
	}

	return nil
}

func deleteAllVersions(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string, bypassGovernance bool) error {
	var keyMarker, versionMarker string
	totalDeleted := 0

	for {
		req := s3.VersionsApi.ListObjectVersions(c.Context, bucket).MaxKeys(1000)
		if keyMarker != "" {
			req = req.KeyMarker(keyMarker)
		}
		if versionMarker != "" {
			req = req.VersionIdMarker(versionMarker)
		}

		result, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("listing object versions: %w", err)
		}

		ids := versionsToIdentifiers(result.GetVersions())
		ids = append(ids, deleteMarkersToIdentifiers(result.GetDeleteMarkers())...)

		if len(ids) == 0 {
			break
		}

		if err := batchDelete(c.Context, s3, bucket, ids, bypassGovernance); err != nil {
			return err
		}

		totalDeleted += len(ids)
		c.Verbose("Deleted %d versions/markers...", totalDeleted)

		if !result.GetIsTruncated() {
			break
		}

		keyMarker = result.GetNextKeyMarker()
		versionMarker = result.GetNextVersionIdMarker()
		if keyMarker == "" {
			break
		}
	}

	return nil
}

func batchDelete(ctx context.Context, s3 *objectstorage.APIClient, bucket string, ids []objectstorage.ObjectIdentifier, bypassGovernance bool) error {
	delReq := objectstorage.NewDeleteObjectsRequest()
	delReq.SetObjects(ids)
	delReq.SetQuiet(true)

	req := s3.ObjectsApi.DeleteObjects(ctx, bucket).DeleteObjectsRequest(*delReq).XAmzBypassGovernanceRetention(bypassGovernance)
	result, _, err := req.Execute()
	if err != nil {
		return fmt.Errorf("deleting objects: %w", err)
	}

	if result != nil && len(result.Errors) > 0 {
		first := result.Errors[0]
		return fmt.Errorf("failed to delete %d object(s): %s: %s (key: %s)",
			len(result.Errors), first.GetCode(), first.GetMessage(), first.GetKey())
	}

	return nil
}

func objectsToIdentifiers(objects []objectstorage.Object) []objectstorage.ObjectIdentifier {
	ids := make([]objectstorage.ObjectIdentifier, len(objects))
	for i, obj := range objects {
		ids[i] = objectstorage.ObjectIdentifier{Key: obj.GetKey()}
	}
	return ids
}

func versionsToIdentifiers(versions []objectstorage.ObjectVersion) []objectstorage.ObjectIdentifier {
	ids := make([]objectstorage.ObjectIdentifier, 0, len(versions))
	for _, v := range versions {
		id := objectstorage.ObjectIdentifier{Key: v.GetKey()}
		if vid := v.GetVersionId(); vid != "" {
			id.VersionId = &vid
		}
		ids = append(ids, id)
	}
	return ids
}

func deleteMarkersToIdentifiers(markers []objectstorage.DeleteMarkerEntry) []objectstorage.ObjectIdentifier {
	ids := make([]objectstorage.ObjectIdentifier, 0, len(markers))
	for _, dm := range markers {
		id := objectstorage.ObjectIdentifier{Key: dm.GetKey()}
		if vid := dm.GetVersionId(); vid != "" {
			id.VersionId = &vid
		}
		ids = append(ids, id)
	}
	return ids
}
