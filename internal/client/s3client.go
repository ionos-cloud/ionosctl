package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// ownerIDFixTransport wraps an http.RoundTripper to rewrite non-numeric
// <ID>...</ID> values inside <Owner> elements to "0". This works around
// the SDK defining Owner.ID as *int32 while the S3 API can return
// "anonymous" as the owner ID.
type ownerIDFixTransport struct {
	base http.RoundTripper
}

var ownerIDRe = regexp.MustCompile(`(<Owner>\s*<ID>)[^<]*(\D)[^<]*(</ID>)`)

func (t *ownerIDFixTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	resp, err := base.RoundTrip(req)
	if err != nil || resp.Body == nil {
		return resp, err
	}

	// Only rewrite XML responses (listing/metadata). Leave object data
	// (downloads, etc.) untouched to preserve streaming and avoid buffering
	// potentially large payloads into memory.
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "xml") {
		return resp, nil
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body = ownerIDRe.ReplaceAll(body, []byte("${1}0${3}"))

	resp.Body = io.NopCloser(bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	return resp, nil
}

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

	opts := shared.ClientOptions{
		Endpoint:            fmt.Sprintf("https://s3.%s.ionoscloud.com", region),
		ObjectStorageRegion: region,
		Credentials: shared.Credentials{
			S3AccessKey: accessKey,
			S3SecretKey: secretKey,
		},
	}
	cfg := shared.NewConfigurationFromOptions(opts).WithObjectStorage(opts)

	// Wrap the transport to fix non-numeric Owner IDs from the S3 API.
	cfg.HTTPClient.Transport = &ownerIDFixTransport{base: cfg.HTTPClient.Transport}

	return objectstorage.NewAPIClient(cfg), nil
}

// GetRegionalObjectStorageClient creates an S3 client targeted at the bucket's
// actual region. It first calls GetBucketLocation via a default-region client,
// then returns a client configured for that region and the resolved region string.
// Use this for any bucket-specific operation that would otherwise hit redirect loops.
func GetRegionalObjectStorageClient(ctx context.Context, bucket string) (*objectstorage.APIClient, string, error) {
	s3, err := GetObjectStorageClient("")
	if err != nil {
		return nil, "", err
	}

	loc, _, err := s3.BucketsApi.GetBucketLocation(ctx, bucket).Execute()
	if err != nil {
		return nil, "", fmt.Errorf("resolving bucket region: %w", err)
	}

	region := ""
	if loc != nil {
		region = loc.GetLocationConstraint()
	}

	cl, err := GetObjectStorageClient(region)
	return cl, region, err
}
