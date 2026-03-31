package bucket

import (
	"context"
	"fmt"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
)

const flagRecursive = "recursive"

func DeleteBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a contract-owned bucket",
		LongDesc:  "Delete a contract-owned bucket. Use --recursive to delete all objects in the bucket first.",
		Example:   "ionosctl object-storage bucket delete --name my-bucket\nionosctl object-storage bucket delete --name my-bucket --recursive -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			recursive := viper.GetBool(core.GetFlagName(c.NS, flagRecursive))

			promptMsg := fmt.Sprintf("delete bucket %q", name)
			if recursive {
				promptMsg = fmt.Sprintf("delete bucket %q and ALL its objects", name)
			}

			if !confirm.FAsk(c.Command.Command.InOrStdin(), promptMsg, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			s3Regional, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			if recursive {
				if err := emptyBucket(c, s3Regional, name); err != nil {
					return err
				}
			}

			_, err = s3Regional.BucketsApi.DeleteBucket(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q deleted successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to delete", core.RequiredFlagOption())
	cmd.AddBoolFlag(flagRecursive, "", false, "Delete all objects in the bucket before deleting the bucket itself")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

// emptyBucket deletes all objects, object versions, and delete markers in the bucket.
func emptyBucket(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string) error {
	// First pass: delete current objects via ListObjectsV2.
	if err := deleteCurrentObjects(c, s3, bucket); err != nil {
		return err
	}

	// Second pass: delete all versions and delete markers via ListObjectVersions.
	if err := deleteAllVersions(c, s3, bucket); err != nil {
		return err
	}

	return nil
}

func deleteCurrentObjects(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string) error {
	var continuationToken string
	totalDeleted := 0

	for {
		req := s3.ObjectsApi.ListObjectsV2(context.Background(), bucket).MaxKeys(1000)
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

		if err := batchDelete(s3, bucket, objectsToIdentifiers(result.Contents)); err != nil {
			return err
		}

		totalDeleted += len(result.Contents)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "Deleted %d objects...\n", totalDeleted)

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

func deleteAllVersions(c *core.CommandConfig, s3 *objectstorage.APIClient, bucket string) error {
	var keyMarker, versionMarker string
	totalDeleted := 0

	for {
		req := s3.VersionsApi.ListObjectVersions(context.Background(), bucket).MaxKeys(1000)
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

		var ids []objectstorage.ObjectIdentifier

		for _, v := range result.GetVersions() {
			id := objectstorage.ObjectIdentifier{Key: v.GetKey()}
			if vid := v.GetVersionId(); vid != "" {
				id.VersionId = &vid
			}
			ids = append(ids, id)
		}

		for _, dm := range result.GetDeleteMarkers() {
			id := objectstorage.ObjectIdentifier{Key: dm.GetKey()}
			if vid := dm.GetVersionId(); vid != "" {
				id.VersionId = &vid
			}
			ids = append(ids, id)
		}

		if len(ids) == 0 {
			break
		}

		if err := batchDelete(s3, bucket, ids); err != nil {
			return err
		}

		totalDeleted += len(ids)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "Deleted %d versions/markers...\n", totalDeleted)

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

func batchDelete(s3 *objectstorage.APIClient, bucket string, ids []objectstorage.ObjectIdentifier) error {
	delReq := objectstorage.NewDeleteObjectsRequest()
	delReq.SetObjects(ids)
	delReq.SetQuiet(true)

	result, _, err := s3.ObjectsApi.DeleteObjects(context.Background(), bucket).
		DeleteObjectsRequest(*delReq).
		Execute()
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
