package bucket

import (
	"context"
	"fmt"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
)

func DeleteBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a contract-owned bucket",
		LongDesc:  "Delete a contract-owned bucket, or all buckets using --all. The bucket must be empty before deletion. Use 'ionosctl object-storage object delete --all' to empty a bucket first.",
		Example:   "ionosctl object-storage bucket delete --name my-bucket\nionosctl object-storage bucket delete --all\nionosctl object-storage bucket delete --all -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.ArgAll},
				[]string{constants.FlagName},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
				return deleteAllBuckets(c)
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete bucket %q", name), viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.DeleteBucket(c.Context, name).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q deleted successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to delete",
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all buckets")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

func deleteAllBuckets(c *core.CommandConfig) error {
	s3 := client.MustObjectStorage().ObjectStorageClient

	result, _, err := s3.BucketsApi.ListBuckets(c.Context).Execute()
	if err != nil {
		return fmt.Errorf("failed listing buckets: %w", err)
	}

	buckets := result.GetBuckets()

	if viper.IsSet(constants.FlagLocation) {
		filterRegion := viper.GetString(constants.FlagLocation)
		filtered, locErr := filterBucketsByLocation(c.Context, s3, buckets, filterRegion)
		if locErr != nil {
			return locErr
		}
		buckets = filtered
	}

	return functional.ApplyAndAggregateErrors(buckets, func(b objectstorage.Bucket) error {
		name := b.GetName()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete bucket %q", name), viper.GetBool(constants.ArgForce)) {
			return nil
		}

		_, delErr := s3.BucketsApi.DeleteBucket(c.Context, name).Execute()
		if delErr != nil {
			return fmt.Errorf("failed deleting bucket %q: %w", name, delErr)
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q deleted successfully\n", name)
		return nil
	})
}

// filterBucketsByLocation filters buckets by region, returning an error if any
// location lookup fails (to avoid silently skipping buckets from deletion).
func filterBucketsByLocation(ctx context.Context, s3 *objectstorage.APIClient, buckets []objectstorage.Bucket, region string) ([]objectstorage.Bucket, error) {
	var filtered []objectstorage.Bucket
	for _, b := range buckets {
		loc, _, locErr := s3.BucketsApi.GetBucketLocation(ctx, b.GetName()).Execute()
		if locErr != nil {
			return nil, fmt.Errorf("cannot filter by location: failed to get location for bucket %q: %w", b.GetName(), locErr)
		}
		if loc.GetLocationConstraint() == region {
			filtered = append(filtered, b)
		}
	}
	return filtered, nil
}
