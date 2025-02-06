package client

import (
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	maria "github.com/ionos-cloud/sdk-go-dbaas-mariadb"
	postgres "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

	"github.com/spf13/viper"
)

var ConfigurationPriorityRules = []Layer{
	{constants.ArgToken, "", "", fmt.Sprintf("Global Flags (--%s)", constants.ArgToken)},
	{
		constants.EnvToken, constants.EnvUsername, constants.EnvPassword,
		fmt.Sprintf(
			"Environment Variables (%s, %s, %s)", constants.EnvToken, constants.EnvUsername, constants.EnvPassword,
		),
	},
	{
		constants.CfgToken, constants.CfgUsername, constants.CfgPassword,
		fmt.Sprintf(
			"Config file settings (%s, %s, %s)", constants.CfgToken, constants.CfgUsername, constants.CfgPassword,
		),
	}, // Note: Username & Password are no longer generated in cfg file by `ionosctl login`, however we will keep this for backward compatibility.
}

// Layer represents an authentication layer. E.g., flags, env vars, config file.
// A client can use one of these layers to authenticate against CloudAPI,
// each layer has priority over layers that are defined after it.
// the Token has priority over username & password pairs of the same authentication layer.
type Layer struct {
	TokenKey    string
	UsernameKey string
	PasswordKey string
	Description string // You can optionally pass a string to describe to the user what this layer is and how to set its values
}

// IsTokenAuth returns true if a token is being used for authentication. Otherwise, username & password were used.
func (c *Client) IsTokenAuth() bool {
	return c.CloudClient.GetConfig().Token != ""
}

func (c *Client) UsedLayer() *Layer {
	if c == nil || c.usedLayer == nil {
		return nil
	}
	return c.usedLayer
}

type Client struct {
	usedLayer *Layer // i.e. which auth layer are we using. Flags / Env Vars / Config File

	CloudClient          *cloudv6.APIClient
	AuthClient           *sdkgoauth.APIClient
	CertManagerClient    *certmanager.APIClient
	DataplatformClient   *dataplatform.APIClient
	RegistryClient       *containerregistry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logging.APIClient
	VMAscClient          *vmasc.AutoScalingGroupsApiService
	VPNClient            *vpn.APIClient

	PostgresClient *postgres.APIClient
	MongoClient    *mongo.APIClient
	MariaClient    *maria.APIClient
	CDNClient      *cdn.APIClient
	Kafka          *kafka.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string, usedLayer *Layer) *Client {
	// TODO: Replace all configurations with this one
	sharedConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	sharedConfig.UserAgent = appendUserAgent(sharedConfig.UserAgent)

	clientConfig := cloudv6.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)
	// Set Depth Query Parameter globally
	clientConfig.SetDepth(1)

	authConfig := sdkgoauth.NewConfiguration(name, pwd, token, hostUrl)
	authConfig.UserAgent = appendUserAgent(authConfig.UserAgent)

	certManagerConfig := certmanager.NewConfiguration(name, pwd, token, hostUrl)
	certManagerConfig.UserAgent = appendUserAgent(certManagerConfig.UserAgent)

	vmascConfig := vmasc.NewConfiguration(name, pwd, token, hostUrl)
	vmascConfig.UserAgent = appendUserAgent(vmascConfig.UserAgent)
	// DBAAS
	postgresConfig := postgres.NewConfiguration(name, pwd, token, hostUrl)
	postgresConfig.UserAgent = appendUserAgent(postgresConfig.UserAgent)

	mariaConfig := maria.NewConfiguration(name, pwd, token, hostUrl)
	mariaConfig.UserAgent = appendUserAgent(mariaConfig.UserAgent)

	return &Client{
		CloudClient:          cloudv6.NewAPIClient(clientConfig),
		AuthClient:           sdkgoauth.NewAPIClient(authConfig),
		CDNClient:            cdn.NewAPIClient(sharedConfig),
		CertManagerClient:    certmanager.NewAPIClient(certManagerConfig),
		RegistryClient:       containerregistry.NewAPIClient(sharedConfig),
		DataplatformClient:   dataplatform.NewAPIClient(sharedConfig),
		DnsClient:            dns.NewAPIClient(sharedConfig),
		LoggingServiceClient: logging.NewAPIClient(sharedConfig),
		VMAscClient:          vmasc.NewAPIClient(vmascConfig).AutoScalingGroupsApi,
		VPNClient:            vpn.NewAPIClient(sharedConfig),

		PostgresClient: postgres.NewAPIClient(postgresConfig),
		MongoClient:    mongo.NewAPIClient(sharedConfig),
		MariaClient:    maria.NewAPIClient(mariaConfig),
		Kafka:          kafka.NewAPIClient(sharedConfig),

		usedLayer: usedLayer,
	}
}
