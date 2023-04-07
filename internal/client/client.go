package client

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	sdkgoauth "github.com/ionos-cloud/sdk-go-bundle/products/auth"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	cloudv6 "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	registry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	dataplatform "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo"
	postgres "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string) (*Client, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("username, password or token incorrect")
	}

	cloudUrl := fmt.Sprintf("%s/cloudapi/v6", hostUrl)
	clientConfig := shared.NewConfiguration(name, pwd, token, cloudUrl)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), clientConfig.UserAgent)
	// Set Depth Query Parameter globally
	clientConfig.DefaultQueryParams.Add("depth", "1")

	authUrl := fmt.Sprintf("%s/auth/v1", hostUrl)
	authConfig := shared.NewConfiguration(name, pwd, token, authUrl)
	authConfig.UserAgent = appendUserAgent(authConfig.UserAgent)

	certManagerConfig := shared.NewConfiguration(name, pwd, token, hostUrl)
	certManagerConfig.UserAgent = appendUserAgent(certManagerConfig.UserAgent)

	postgresUrl := fmt.Sprintf("%s/databases/postgresql", hostUrl)
	postgresConfig := shared.NewConfiguration(name, pwd, token, postgresUrl)
	postgresConfig.UserAgent = appendUserAgent(postgresConfig.UserAgent)

	mongoUrl := fmt.Sprintf("%s/databases/mongodb", hostUrl)
	mongoConfig := shared.NewConfiguration(name, pwd, token, mongoUrl)
	mongoConfig.UserAgent = appendUserAgent(mongoConfig.UserAgent)

	dataplatformURL := fmt.Sprintf("%s/dataplatform", hostUrl)
	dpConfig := shared.NewConfiguration(name, pwd, token, dataplatformURL)
	dpConfig.UserAgent = appendUserAgent(dpConfig.UserAgent)

	contregURL := fmt.Sprintf("%s/containerregistries", hostUrl)
	registryConfig := shared.NewConfiguration(name, pwd, token, contregURL)
	registryConfig.UserAgent = appendUserAgent(registryConfig.UserAgent)

	return &Client{
			CloudClient:        cloudv6.NewAPIClient(clientConfig),
			AuthClient:         sdkgoauth.NewAPIClient(authConfig),
			CertManagerClient:  certmanager.NewAPIClient(certManagerConfig),
			PostgresClient:     postgres.NewAPIClient(postgresConfig),
			MongoClient:        mongo.NewAPIClient(mongoConfig),
			DataplatformClient: dataplatform.NewAPIClient(dpConfig),
			RegistryClient:     registry.NewAPIClient(registryConfig),
		},
		nil
}

var once sync.Once
var instance *Client

// Get a client and possibly fail
func Get() (*Client, error) {
	var err error
	once.Do(func() {
		err = config.Load()
		instance, err = newClient(viper.GetString(constants.Username), viper.GetString(constants.Password), viper.GetString(constants.Token), config.GetServerUrl())
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
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
