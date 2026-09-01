package object

import (
	"context"
	"fmt"
	"strings"

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
		ShortDesc: "Delete a single object, or empty a whole bucket",
		LongDesc: `Delete one object by key, or empty the entire bucket with --all.

Versioning changes what "delete" means:
  - On a versioning-enabled bucket, deleting a key WITHOUT --version-id does not remove data: it inserts a "delete marker" that hides the key, while all prior versions remain recoverable.
  - Passing --version-id permanently removes that one specific version (this cannot be undone). Deleting a delete-marker version un-hides the key.
  - On a bucket without versioning, the object is simply removed.

--all empties the bucket: it deletes every current object AND every historical version AND every delete marker, so the bucket ends up truly empty. This is destructive and irreversible - it is guarded by a confirmation prompt (use -f/--force to skip it in scripts).

Object Lock: objects protected by a retention or legal hold cannot be deleted. GOVERNANCE-mode retention can be overridden with --bypass-governance-retention (requires the appropriate permission); COMPLIANCE-mode retention and active legal holds can never be bypassed.`,
		Example: `# Delete one object (inserts a delete marker on a versioned bucket)
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg

# Permanently delete one specific version
ionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id <version-id> -f

# Empty the whole bucket (all objects, versions and delete markers)
ionosctl object-storage object delete --name my-bucket --all -f

# Empty a bucket, overriding GOVERNANCE-mode Object Lock protection
ionosctl object-storage object delete --name my-bucket --all --bypass-governance-retention -f`,
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to delete from", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object to delete. Mutually exclusive with --all",
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Permanently delete this specific version (irreversible). Without it, deleting on a versioned bucket only inserts a delete marker")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Empty the entire bucket: delete every object, version and delete marker. Destructive and irreversible")
	cmd.AddBoolFlag(flagBypassGovernanceRetention, "", false, "Override GOVERNANCE-mode Object Lock retention so locked objects can be deleted (needs bypass permission). Has no effect on COMPLIANCE-mode retention or legal holds")

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
		var details []string
		for _, e := range result.Errors {
			details = append(details, fmt.Sprintf("  %s: %s (key: %s)", e.GetCode(), e.GetMessage(), e.GetKey()))
		}
		return fmt.Errorf("failed to delete %d object(s):\n%s",
			len(result.Errors), strings.Join(details, "\n"))
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
