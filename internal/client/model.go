package client

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	registry "github.com/ionos-cloud/sdk-go-container-registry"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	postgres "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	dns "github.com/ionos-cloud/sdk-go-dns"
	logsvc "github.com/ionos-cloud/sdk-go-logging"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

var ConfigurationPriorityRules = []Layer{
	{constants.ArgToken, "", "", fmt.Sprintf("Global Flags (--%s)", constants.ArgToken)},
	{
		constants.EnvToken, constants.EnvUsername, constants.EnvPassword, fmt.Sprintf(
			"Environment Variables (%s, %s, %s)", constants.EnvToken, constants.EnvUsername, constants.EnvPassword,
		),
	},
	{
		constants.CfgToken, constants.CfgUsername, constants.CfgPassword, fmt.Sprintf(
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
	PostgresClient       *postgres.APIClient
	MongoClient          *mongo.APIClient
	DataplatformClient   *dataplatform.APIClient
	RegistryClient       *registry.APIClient
	DnsClient            *dns.APIClient
	LoggingServiceClient *logsvc.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string, usedLayer *Layer) *Client {
	clientConfig := cloudv6.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = appendUserAgent(clientConfig.UserAgent)
	// Set Depth Query Parameter globally
	clientConfig.SetDepth(1)

	authConfig := sdkgoauth.NewConfiguration(name, pwd, token, hostUrl)
	authConfig.UserAgent = appendUserAgent(authConfig.UserAgent)

	certManagerConfig := certmanager.NewConfiguration(name, pwd, token, hostUrl)
	certManagerConfig.UserAgent = appendUserAgent(certManagerConfig.UserAgent)

	postgresConfig := postgres.NewConfiguration(name, pwd, token, hostUrl)
	postgresConfig.UserAgent = appendUserAgent(postgresConfig.UserAgent)

	mongoConfig := mongo.NewConfiguration(name, pwd, token, hostUrl)
	mongoConfig.UserAgent = appendUserAgent(mongoConfig.UserAgent)

	dpConfig := dataplatform.NewConfiguration(name, pwd, token, hostUrl)
	dpConfig.UserAgent = appendUserAgent(dpConfig.UserAgent)

	registryConfig := registry.NewConfiguration(name, pwd, token, hostUrl)
	registryConfig.UserAgent = appendUserAgent(registryConfig.UserAgent)

	dnsConfig := dns.NewConfiguration(name, pwd, token, hostUrl)
	dnsConfig.UserAgent = appendUserAgent(dnsConfig.UserAgent)

	logsConfig := logsvc.NewConfiguration(name, pwd, token, hostUrl)
	logsConfig.UserAgent = appendUserAgent(logsConfig.UserAgent)

	return &Client{
		CloudClient:          cloudv6.NewAPIClient(clientConfig),
		AuthClient:           sdkgoauth.NewAPIClient(authConfig),
		CertManagerClient:    certmanager.NewAPIClient(certManagerConfig),
		PostgresClient:       postgres.NewAPIClient(postgresConfig),
		MongoClient:          mongo.NewAPIClient(mongoConfig),
		DataplatformClient:   dataplatform.NewAPIClient(dpConfig),
		RegistryClient:       registry.NewAPIClient(registryConfig),
		DnsClient:            dns.NewAPIClient(dnsConfig),
		LoggingServiceClient: logsvc.NewAPIClient(logsConfig),
		usedLayer:            usedLayer,
	}
}
