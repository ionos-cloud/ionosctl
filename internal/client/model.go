package client

import (
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	postgres "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

	"github.com/spf13/viper"
)

// AuthSource represents a human-readable description of where the client's authentication credentials were sourced from.
type AuthSource string

const (
	AuthSourceEnvBearer AuthSource = "environment variable: IONOS_TOKEN"
	AuthSourceEnvBasic  AuthSource = "environment variables: IONOS_USERNAME, IONOS_PASSWORD"
	AuthSourceCfgBearer AuthSource = "token from config file: ."
	AuthSourceNone      AuthSource = "no authentication provided"
)

type Client struct {
	Config     *fileconfiguration.FileConfig
	AuthSource AuthSource

	Apigateway           *apigateway.APIClient
	CloudClient          *cloudv6.APIClient
	AuthClient           *auth.APIClient
	CertManagerClient    *certmanager.APIClient
	DataplatformClient   *dataplatform.APIClient
	RegistryClient       *containerregistry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logging.APIClient
	VMAscClient          *vmasc.AutoScalingGroupsApiService
	VPNClient            *vpn.APIClient
	CDNClient            *cdn.APIClient
	Kafka                *kafka.APIClient
	Monitoring           *monitoring.APIClient

	PostgresClient   *postgres.APIClient
	MongoClient      *mongo.APIClient
	MariaClient      *mariadb.APIClient
	InMemoryDBClient *inmemorydb.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string) *Client {
	// TODO: Replace all configurations with this one
	sharedConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	sharedConfig.UserAgent = appendUserAgent(sharedConfig.UserAgent)

	clientConfig := cloudv6.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)
	// Set Depth Query Parameter globally
	clientConfig.SetDepth(1)

	certManagerConfig := certmanager.NewConfiguration(name, pwd, token, hostUrl)
	certManagerConfig.UserAgent = appendUserAgent(certManagerConfig.UserAgent)

	vmascConfig := vmasc.NewConfiguration(name, pwd, token, hostUrl)
	vmascConfig.UserAgent = appendUserAgent(vmascConfig.UserAgent)
	// DBAAS
	postgresConfig := postgres.NewConfiguration(name, pwd, token, hostUrl)
	postgresConfig.UserAgent = appendUserAgent(postgresConfig.UserAgent)

	return &Client{
		URLOverride: hostUrl,

		Apigateway:           apigateway.NewAPIClient(sharedConfig),
		CloudClient:          cloudv6.NewAPIClient(clientConfig),
		AuthClient:           auth.NewAPIClient(sharedConfig),
		CDNClient:            cdn.NewAPIClient(sharedConfig),
		CertManagerClient:    certmanager.NewAPIClient(certManagerConfig),
		RegistryClient:       containerregistry.NewAPIClient(sharedConfig),
		DataplatformClient:   dataplatform.NewAPIClient(sharedConfig),
		DnsClient:            dns.NewAPIClient(sharedConfig),
		LoggingServiceClient: logging.NewAPIClient(sharedConfig),
		VMAscClient:          vmasc.NewAPIClient(vmascConfig).AutoScalingGroupsApi,
		VPNClient:            vpn.NewAPIClient(sharedConfig),
		Kafka:                kafka.NewAPIClient(sharedConfig),
		Monitoring:           monitoring.NewAPIClient(sharedConfig),

		PostgresClient:   postgres.NewAPIClient(postgresConfig),
		MongoClient:      mongo.NewAPIClient(sharedConfig),
		MariaClient:      mariadb.NewAPIClient(sharedConfig),
		InMemoryDBClient: inmemorydb.NewAPIClient(sharedConfig),
	}
}
