package client

import (
	"fmt"
	"os"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetObjectStorageClient returns an S3-authenticated APIClient for the given region.
// Credentials are resolved in priority order:
//  1. Environment variables IONOS_S3_ACCESS_KEY / IONOS_S3_SECRET_KEY
//  2. s3AccessKey / s3SecretKey in the current ionosctl config profile
func GetObjectStorageClient(region string) (*objectstorage.APIClient, error) {
	accessKey := os.Getenv(shared.IonosS3AccessKeyEnvVar)
	secretKey := os.Getenv(shared.IonosS3SecretKeyEnvVar)

	if accessKey == "" || secretKey == "" {
		if c, err := Get(); err == nil && c.Config != nil && c.Config.GetCurrentProfile() != nil {
			creds := c.Config.GetCurrentProfile().Credentials
			if accessKey == "" {
				accessKey = creds.S3AccessKey
			}
			if secretKey == "" {
				secretKey = creds.S3SecretKey
			}
		}
	}

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf(
			"object storage credentials not found. Set %s and %s environment variables, or configure s3AccessKey/s3SecretKey in your ionosctl profile",
			shared.IonosS3AccessKeyEnvVar, shared.IonosS3SecretKeyEnvVar,
		)
	}

	if region == "" {
		region = "eu-central-3"
	}

	cfg := shared.NewConfigurationFromOptions(shared.ClientOptions{
		Endpoint:            fmt.Sprintf("https://s3.%s.ionoscloud.com", region),
		ObjectStorageRegion: region,
		Credentials: shared.Credentials{
			S3AccessKey: accessKey,
			S3SecretKey: secretKey,
		},
	})

	return objectstorage.NewAPIClient(cfg), nil
}
