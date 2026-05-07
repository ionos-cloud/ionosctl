package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
)

// ownerIDFixTransport wraps an http.RoundTripper to rewrite non-numeric
// <ID>...</ID> values inside <Owner> elements to "0". This works around
// the SDK defining Owner.ID as *int32 while the Object Storage API can return
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

var (
	osOnce      sync.Once
	osInstance  *ObjectStorageClient
	osClientErr error
)

// ResolveObjectStorageCredentials resolves Object Storage access and secret keys, tracking their source.
// Both keys must come from the same source; mixing is not allowed. Priority order:
//  1. Environment variables IONOS_S3_ACCESS_KEY / IONOS_S3_SECRET_KEY (both must be set)
//  2. s3AccessKey / s3SecretKey in the current ionosctl config profile (both must be set)
func ResolveObjectStorageCredentials() (accessKey, secretKey string, akSrc ObjectStorageAccessKeySource, skSrc ObjectStorageSecretKeySource, err error) {
	// Priority 1: both from environment variables
	envAccessKey := os.Getenv(shared.IonosS3AccessKeyEnvVar)
	envSecretKey := os.Getenv(shared.IonosS3SecretKeyEnvVar)
	if envAccessKey != "" && envSecretKey != "" {
		return envAccessKey, envSecretKey, ObjectStorageAccessKeyEnv, ObjectStorageSecretKeyEnv, nil
	}

	// Priority 2: both from config file
	src, cfgErr := retrieveConfigFile()
	if cfgErr != nil {
		return "", "", ObjectStorageAccessKeyNone, ObjectStorageSecretKeyNone,
			fmt.Errorf("failed to retrieve config file: %w", cfgErr)
	}
	if src.Config != nil && src.Config.GetCurrentProfile() != nil {
		creds := src.Config.GetCurrentProfile().Credentials
		if creds.S3AccessKey != "" && creds.S3SecretKey != "" {
			return creds.S3AccessKey, creds.S3SecretKey, ObjectStorageAccessKeyCfg, ObjectStorageSecretKeyCfg, nil
		}
	}

	return "", "", ObjectStorageAccessKeyNone, ObjectStorageSecretKeyNone, fmt.Errorf(
		"object storage credentials not found. Set %s and %s environment variables, or configure s3AccessKey/s3SecretKey in your ionosctl profile",
		shared.IonosS3AccessKeyEnvVar, shared.IonosS3SecretKeyEnvVar,
	)
}

// newObjectStorageClient builds a new ObjectStorageClient for the given endpoint.
func newObjectStorageClient(endpoint, region string) (*ObjectStorageClient, error) {
	accessKey, secretKey, akSrc, skSrc, err := ResolveObjectStorageCredentials()
	if err != nil {
		return nil, err
	}

	opts := shared.ClientOptions{
		Endpoint:            endpoint,
		ObjectStorageRegion: region,
		Credentials: shared.Credentials{
			S3AccessKey: accessKey,
			S3SecretKey: secretKey,
		},
	}
	cfg := shared.NewConfigurationFromOptions(opts).WithObjectStorage(opts)

	// Wrap the transport to fix non-numeric Owner IDs from the Object Storage API.
	cfg.HTTPClient.Transport = &ownerIDFixTransport{base: cfg.HTTPClient.Transport}

	return &ObjectStorageClient{
		AccessKeySource:     akSrc,
		SecretKeySource:     skSrc,
		URLOverride:         endpoint,
		Region:              region,
		ObjectStorageClient: objectstorage.NewAPIClient(cfg),
	}, nil
}

// GetObjectStorage returns the ObjectStorageClient for the currently resolved endpoint.
// The endpoint is set by WithRegionalConfigOverride via viper (constants.ArgServerUrl).
// Falls back to the default region endpoint if not set. Cached via sync.Once.
func GetObjectStorage() (*ObjectStorageClient, error) {
	osOnce.Do(func() {
		endpoint := viper.GetString(constants.ArgServerUrl)
		region := viper.GetString(constants.FlagLocation)

		if endpoint == "" {
			if region == "" {
				region = constants.ObjectStorageLocations[0]
			}
			endpoint = fmt.Sprintf(constants.ObjectStorageApiRegionalURL, region)
		}

		osInstance, osClientErr = newObjectStorageClient(endpoint, region)
	})

	return osInstance, osClientErr
}

var MustObjectStorageDefaultErrHandler = func(err error) {
	die.Die(fmt.Errorf("failed getting object storage client: %w", err).Error())
}

// MustObjectStorage returns the ObjectStorageClient or fatally exits.
// Custom error handlers can be provided; the default handler calls die.Die.
func MustObjectStorage(ehs ...func(error)) *ObjectStorageClient {
	cl, err := GetObjectStorage()
	if err != nil {
		if len(ehs) > 0 {
			for _, eh := range ehs {
				eh(err)
			}
		} else {
			MustObjectStorageDefaultErrHandler(err)
		}
	}
	return cl
}
