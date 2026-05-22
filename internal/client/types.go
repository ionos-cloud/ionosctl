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

type ObjectStorageAccessKeySource string
type ObjectStorageSecretKeySource string

const (
	ObjectStorageAccessKeyEnv  ObjectStorageAccessKeySource = "environment variable: IONOS_S3_ACCESS_KEY"
	ObjectStorageAccessKeyCfg  ObjectStorageAccessKeySource = "Object Storage access key from config file: s3AccessKey"
	ObjectStorageAccessKeyNone ObjectStorageAccessKeySource = "Object Storage access key not provided"

	ObjectStorageSecretKeyEnv  ObjectStorageSecretKeySource = "environment variable: IONOS_S3_SECRET_KEY"
	ObjectStorageSecretKeyCfg  ObjectStorageSecretKeySource = "Object Storage secret key from config file: s3SecretKey"
	ObjectStorageSecretKeyNone ObjectStorageSecretKeySource = "Object Storage secret key not provided"
)

// all possible sources in priority order
var AuthOrder = []AuthSource{
	AuthSourceEnvBearer,
	AuthSourceEnvBasic,
	AuthSourceCfgBearer,
	AuthSourceCfgBasic,
}

var ObjectStorageAccessKeyOrder = []ObjectStorageAccessKeySource{
	ObjectStorageAccessKeyEnv,
	ObjectStorageAccessKeyCfg,
}

var ObjectStorageSecretKeyOrder = []ObjectStorageSecretKeySource{
	ObjectStorageSecretKeyEnv,
	ObjectStorageSecretKeyCfg,
}

type Client struct {
	Config     *fileconfiguration.FileConfig
	ConfigPath string // Path to the config file used to create this client, if any.
	AuthSource AuthSource

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
	AccessKeySource ObjectStorageAccessKeySource
	SecretKeySource ObjectStorageSecretKeySource

	URLOverride string
	Region      string

	ObjectStorageClient *objectstorage.APIClient
}
