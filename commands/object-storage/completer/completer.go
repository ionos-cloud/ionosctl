package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

// BucketNames returns all bucket names for shell autocompletion.
func BucketNames() []string {
	s3, err := client.GetObjectStorageClient("")
	if err != nil {
		return nil
	}

	result, _, err := s3.BucketsApi.ListBuckets(context.Background()).Execute()
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

	s3, _, err := client.GetRegionalObjectStorageClient(context.Background(), bucket)
	if err != nil {
		return nil
	}

	result, _, err := s3.ObjectsApi.ListObjectsV2(context.Background(), bucket).
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

// Regions returns available S3 regions.
func Regions() []string {
	return []string{"eu-central-3", "eu-central-4", "us-central-1"}
}
