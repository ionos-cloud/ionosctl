package client

import (
	"context"
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
}

func appendUserAgent(userAgent string) string {
	return fmt.Sprintf("%v_%v", viper.GetString(constants.CLIHttpUserAgent), userAgent)
}

func newClient(name, pwd, token, hostUrl string) (*Client, error) {
	if token == "" && (name == "" || pwd == "") {
		return nil, errors.New("both token and at least one of username and password are empty")
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

// Get a client and possibly fail. Uses viper to get the credentials and API URL.
// NOTE: The Viper credentials are bound using BindEnv, in this order:
// `--token`, `IONOS_TOKEN`, config file's `userdata.token`
// `--api-url`, `IONOS_API_URL`, config file api url.
// `IONOS_USERNAME`, config file's `userdata.user` (Config file entry is kept for backwards compatibility reasons. However, the config file
func Get() (*Client, error) {
	var getClientErr error

	once.Do(func() {
		var err error

		data, err := config.Read()
		if err != nil {
			// Failed reading config
			if viper.IsSet(constants.ArgToken) || (viper.IsSet(constants.EnvUsername) && viper.IsSet(constants.EnvPassword)) {
				// It's fine if we got the credentials from some place else though (eg env vars)
				err = testCredentialsFromViper()
				if err != nil {
					getClientErr = errors.Join(getClientErr, err)
					return
				}
			}
			getClientErr = errors.Join(getClientErr, fmt.Errorf("failed reading config file: %w", err))
			return
		}

		for k, v := range data {
			// Load config data into viper if not set
			if !viper.IsSet(k) {
				viper.Set(k, v)
			}
		}

		if !viper.IsSet(constants.ArgToken) && (!viper.IsSet(constants.EnvUsername) || !viper.IsSet(constants.EnvPassword)) {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("not logged in: use either environment variables %s or %s and %s, either `ionosctl login`", constants.EnvToken, constants.EnvUsername, constants.EnvPassword))
			return
		}

		err = testCredentialsFromViper()
		if err != nil {
			getClientErr = errors.Join(getClientErr, err)
			return
		}

		instance, err = newClient(viper.GetString(constants.EnvUsername), viper.GetString(constants.EnvPassword), viper.GetString(constants.ArgToken), viper.GetString(constants.ArgServerUrl))
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

// NewClient bypasses the singleton check, not recommended for normal use.
// Use it if you must
func NewClient(name, pwd, token, hostUrl string) (*Client, error) {
	return newClient(name, pwd, token, hostUrl)
}

func TestCreds(user, pass, token string) error {
	cl, err := newClient(user, pass, token, config.GetServerUrl())
	if err != nil {
		return fmt.Errorf("failed initializing client with credentials: %w", err)
	}

	_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		usedScheme := "used token"
		if token == "" {
			usedScheme = fmt.Sprintf("used username '%s' and password", user)
		}
		return fmt.Errorf("credentials test failed. %s: %w", usedScheme, err)
	}

	return nil
}

// helper used in Get
func testCredentialsFromViper() error {
	return TestCreds(viper.GetString(constants.EnvUsername), viper.GetString(constants.EnvPassword), viper.GetString(constants.ArgToken))
}
