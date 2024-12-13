package client

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	sdkgoauth "github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	cdn "github.com/ionos-cloud/sdk-go-bundle/products/cdn/v2"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	cloudv6 "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	registry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	postgres "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	logsvc "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	maria "github.com/ionos-cloud/sdk-go-dbaas-mariadb"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"

	"github.com/spf13/viper"
)

var ConfigurationPriorityRules = []Layer{
	{constants.ArgToken, "", "", fmt.Sprintf("Global Flags (--%s)", constants.ArgToken)},
	{
		constants.EnvToken, constants.EnvUsername, constants.EnvPassword,
		fmt.Sprintf("Environment Variables (%s, %s, %s)", constants.EnvToken, constants.EnvUsername, constants.EnvPassword),
	},
	{
		constants.CfgToken, constants.CfgUsername, constants.CfgPassword,
		fmt.Sprintf("Config file settings (%s, %s, %s)", constants.CfgToken, constants.CfgUsername, constants.CfgPassword),
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
	RegistryClient       *registry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logsvc.APIClient
	VMAscClient          *vmasc.AutoScalingGroupsApiService

	PostgresClient *postgres.APIClient
	MongoClient    *mongo.APIClient
	MariaClient    *maria.APIClient
	CDNClient      *cdn.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string, usedLayer *Layer) *Client {
	clientConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)

	mariaConfig := maria.NewConfiguration(name, pwd, token, hostUrl)
	mariaConfig.UserAgent = appendUserAgent(mariaConfig.UserAgent)

	vmascConfig := vmasc.NewConfiguration(name, pwd, token, hostUrl)
	vmascConfig.UserAgent = appendUserAgent(vmascConfig.UserAgent)

	return &Client{
		CloudClient:          cloudv6.NewAPIClient(clientConfig),
		AuthClient:           sdkgoauth.NewAPIClient(clientConfig),
		CertManagerClient:    certmanager.NewAPIClient(clientConfig),
		DataplatformClient:   dataplatform.NewAPIClient(clientConfig),
		RegistryClient:       registry.NewAPIClient(clientConfig),
		DnsClient:            dns.NewAPIClient(clientConfig),
		LoggingServiceClient: logsvc.NewAPIClient(clientConfig),
		VMAscClient:          vmasc.NewAPIClient(vmascConfig).AutoScalingGroupsApi,

		PostgresClient: postgres.NewAPIClient(clientConfig),
		MongoClient:    mongo.NewAPIClient(clientConfig),
		MariaClient:    maria.NewAPIClient(mariaConfig),
		CDNClient:      cdn.NewAPIClient(clientConfig),

		usedLayer: usedLayer,
	}
}
