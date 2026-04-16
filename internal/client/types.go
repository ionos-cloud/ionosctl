package client

import (
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	psql2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
)

// AuthSource represents a human-readable description of where the client's authentication credentials were sourced from.
type AuthSource string

const (
	AuthSourceEnvBearer AuthSource = "environment variable: IONOS_TOKEN"
	AuthSourceEnvBasic  AuthSource = "environment variables: IONOS_USERNAME, IONOS_PASSWORD"
	AuthSourceCfgBearer AuthSource = "credentials from config file: token"
	AuthSourceCfgBasic  AuthSource = "credentials from config file: username, password"
	AuthSourceNone      AuthSource = "no authentication provided"
)

type S3AccessKeySource string
type S3SecretKeySource string

const (
	S3AccessKeyEnv  S3AccessKeySource = "environment variable: IONOS_S3_ACCESS_KEY"
	S3AccessKeyCfg  S3AccessKeySource = "S3 access key from config file: s3AccessKey"
	S3AccessKeyNone S3AccessKeySource = "S3 access key not provided"

	S3SecretKeyEnv  S3SecretKeySource = "environment variable: IONOS_S3_SECRET_KEY"
	S3SecretKeyCfg  S3SecretKeySource = "S3 secret key from config file: s3SecretKey"
	S3SecretKeyNone S3SecretKeySource = "S3 secret key not provided"
)

// all possible sources in priority order
var AuthOrder = []AuthSource{
	AuthSourceEnvBearer,
	AuthSourceEnvBasic,
	AuthSourceCfgBearer,
	AuthSourceCfgBasic,
}

type Client struct {
	Config     *fileconfiguration.FileConfig
	ConfigPath string // Path to the config file used to create this client, if any.
	AuthSource AuthSource

	URLOverride string // If the client was created with a specific URL override, this will hold that value. If we notice a change in the URL, we need to re-create the client.

	CloudClient          *cloudv6.APIClient
	AuthClient           *auth.APIClient
	CertManagerClient    *cert.APIClient
	RegistryClient       *containerregistry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logging.APIClient
	VMAscClient          *vmasc.AutoScalingGroupsApiService
	VPNClient            *vpn.APIClient
	CDNClient            *cdn.APIClient
	Kafka                *kafka.APIClient
	Monitoring           *monitoring.APIClient

	PostgresClient   *psql.APIClient
	PostgresClientV2 *psql2.APIClient
	MongoClient      *mongo.APIClient
	MariaClient      *mariadb.APIClient
	InMemoryDBClient *inmemorydb.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

type ObjectStorageClient struct {
	S3SecretAccessKeySource S3AccessKeySource
	S3SecretKeySource       S3SecretKeySource

	URLOverride string
	Region      string

	ObjectStorageClient *objectstorage.APIClient
}
