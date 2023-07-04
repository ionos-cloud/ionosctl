package client

import (
	"errors"
	"fmt"
	"strings"
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
	dns "github.com/ionos-cloud/sdk-go-dns"
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

	// Hacky workaround for not being able to change defaults of root level flags, in sub-commands.
	// Rationale: If it is the global default, it probably didn't get changed by the user
	// FIXME: Replace me with a better alternative. What if the user specifically sets api.ionos.com for DNS ? (Presently, he should get 404, but not sure about future)
	dnsHostUrl := hostUrl // Absolutely do not use `dnsHostUrl` anywhere other than DNS Configuration
	if strings.Trim("https://", hostUrl) == strings.Trim("https://", constants.DefaultApiURL) {
		dnsHostUrl = ""
	}
	dnsConfig := dns.NewConfiguration(name, pwd, token, dnsHostUrl)
	dnsConfig.UserAgent = appendUserAgent(dnsConfig.UserAgent)

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
