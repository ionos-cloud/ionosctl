package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
)

// BucketNames returns all bucket names for shell autocompletion.
func BucketNames() []string {
	osClient, err := client.GetObjectStorage()
	if err != nil {
		return nil
	}

	result, _, err := osClient.ObjectStorageClient.BucketsApi.ListBuckets(context.Background()).Execute()
	if err != nil {
		return nil
	}

	var names []string
	for _, b := range result.GetBuckets() {
		names = append(names, b.GetName())
	}
	return names
}

// ObjectKeys returns object keys for a given bucket.
func ObjectKeys(bucket string) []string {
	if bucket == "" {
		return nil
	}

	osClient, err := client.GetObjectStorage()
	if err != nil {
		return nil
	}

	result, _, err := osClient.ObjectStorageClient.ObjectsApi.ListObjectsV2(context.Background(), bucket).
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

// Regions returns available object storage regions.
func Regions() []string {
	return constants.ObjectStorageLocations
}
