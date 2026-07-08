package completer

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
)

// BucketNames returns all bucket names for shell autocompletion.
func BucketNames() []string {
	result, _, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.ListBuckets(context.Background()).Execute()
	if err != nil {
		return nil
	}

	// Only filter by location when --location was explicitly set; otherwise
	// list buckets from every region (matching `bucket list`). Filtering
	// against an unset location would drop every bucket.
	filterByLocation := viper.IsSet(constants.FlagLocation)
	wantLocation := viper.GetString(constants.FlagLocation)

	var names []string
	for _, b := range result.GetBuckets() {
		if filterByLocation {
			loc, _, locErr := client.MustObjectStorage().ObjectStorageClient.BucketsApi.GetBucketLocation(context.Background(), b.GetName()).Execute()
			if locErr != nil || loc.GetLocationConstraint() != wantLocation {
				continue
			}
		}

		names = append(names, b.GetName())
	}
	return names
}

// ObjectKeys returns object keys for a given bucket.
func ObjectKeys(bucket string) []string {
	if bucket == "" {
		return nil
	}

	result, _, err := client.MustObjectStorage().ObjectStorageClient.ObjectsApi.ListObjectsV2(context.Background(), bucket).
		MaxKeys(100).
		Execute()
	if err != nil {
		return nil
	}

	var keys []string
	for _, obj := range result.Contents {
		keys = append(keys, obj.GetKey())
	}
	return keys
}
