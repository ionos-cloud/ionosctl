package client

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	registry "github.com/ionos-cloud/sdk-go-container-registry"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	mongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	postgres "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	dns "github.com/ionos-cloud/sdk-go-dnsaas"
	cloudv6 "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

type Client struct {
	CloudClient        *cloudv6.APIClient
	AuthClient         *sdkgoauth.APIClient
	CertManagerClient  *certmanager.APIClient
	PostgresClient     *postgres.APIClient
	MongoClient        *mongo.APIClient
	DataplatformClient *dataplatform.APIClient
	RegistryClient     *registry.APIClient
	DnsClient          *dns.APIClient
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string) (*Client, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}

	clientConfig := cloudv6.NewConfiguration(name, pwd, token, hostUrl)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), clientConfig.UserAgent)
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

	// HostURL is disabled for DNS currently because of multiple reasons.
	// 1. https://hosting-jira.1and1.org/browse/SDK-1452
	// 2. I can't override defaults for this flag on a leaf-command-level because config.GetServerURL() then prefers some other default
	// 3. I can't use a decorator cmdRun approach, where I define a custom behaviour for hostUrl at dns namespace level, because of circular imports.
	// I think the easiest way to add back hostUrl is by 2. and fixing GetServerURL
	dnsConfig := dns.NewConfiguration(name, pwd, token, "")
	dnsConfig.UserAgent = appendUserAgent(authConfig.UserAgent)

	return &Client{
			CloudClient:        cloudv6.NewAPIClient(clientConfig),
			AuthClient:         sdkgoauth.NewAPIClient(authConfig),
			CertManagerClient:  certmanager.NewAPIClient(certManagerConfig),
			PostgresClient:     postgres.NewAPIClient(postgresConfig),
			MongoClient:        mongo.NewAPIClient(mongoConfig),
			DataplatformClient: dataplatform.NewAPIClient(dpConfig),
			RegistryClient:     registry.NewAPIClient(registryConfig),
			DnsClient:          dns.NewAPIClient(dnsConfig),
		},
		nil
}

var once sync.Once
var instance *Client

// Get a client and possibly fail
func Get() (*Client, error) {
	var getClientErr error

	once.Do(func() {
		var err error
		err = config.Load()
		if err != nil {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("failed loading config: %w", err))
		}
		instance, err = newClient(viper.GetString(constants.Username), viper.GetString(constants.Password), viper.GetString(constants.Token), config.GetServerUrl())
		if err != nil {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("failed creating client: %w", err))
		}
	})

	return instance, getClientErr
}

// Must gets the client obj or fatally dies
func Must() *Client {
	client, err := Get()
	if err != nil {
		die.Die(fmt.Errorf("failed getting client: %w", err).Error())
	}
	return client
}

// NewTestClient - function used only for tests.
// Bypasses the singleton check, not recommended for normal use.
// TO BE REMOVED ONCE TESTS ARE REFACTORED
func NewTestClient(name, pwd, token, hostUrl string) (*Client, error) {
	return newClient(name, pwd, token, hostUrl)
}
