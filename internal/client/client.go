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
	UsedLayer Layer // i.e. which auth layer are we using. Flags / Env Vars / Config File

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

func newClient(name, pwd, token, hostUrl string) *Client {
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
	}
}

var once sync.Once
var instance *Client

func getFirstValidSource(layers []Layer) (values map[string]string, usedLayer Layer) {
	for _, layer := range layers {
		token := viper.GetString(layer.TokenKey)
		username := viper.GetString(layer.UsernameKey)
		password := viper.GetString(layer.PasswordKey)

		if layer.TokenKey != "" && token != "" ||
			layer.UsernameKey != "" && username != "" && layer.PasswordKey != "" && password != "" {
			return map[string]string{
				"token":    token,
				"username": username,
				"password": password,
			}, layer
		}
	}
	return nil, Layer{}
}

type Layer struct {
	TokenKey    string
	UsernameKey string
	PasswordKey string
	Help        string
}

var ConfigurationPriorityRules = []Layer{
	{constants.ArgToken, "", "", fmt.Sprintf("Global Flags (--%s)", constants.ArgToken)},
	{constants.EnvToken, constants.EnvUsername, constants.EnvPassword, fmt.Sprintf("Environment Variables (%s, %s, %s)", constants.EnvToken, constants.EnvUsername, constants.EnvPassword)},
	{constants.CfgToken, constants.CfgUsername, constants.CfgPassword, fmt.Sprintf("Config file settings (%s, %s, %s)", constants.CfgToken, constants.CfgUsername, constants.CfgPassword)}, // Note: Username & Password are no longer generated in cfg file by `ionosctl login`, however we will keep this for backward compatibility.
}

// Get a client and possibly fail. Uses viper to get the credentials and API URL.
// The returned client is guaranteed to have working credentials
// Order:
// Explicit flags ( e.g. --token )
// Environment Variables ( e.g. IONOS_TOKEN )
// Config File ( e.g. userdata.token )
func Get() (*Client, error) {
	var getClientErr error

	once.Do(func() {
		var err error

		// Read config file, if available
		data, err := config.Read()
		if err == nil {
			for k, v := range data {
				if !viper.IsSet(k) {
					viper.Set(k, v)
				}
			}
		}

		viper.AutomaticEnv()

		values, usedLayer := getFirstValidSource(ConfigurationPriorityRules)
		if values == nil {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("not logged in: use either environment variables %s or %s and %s, or use `ionosctl login`", constants.EnvToken, constants.EnvUsername, constants.EnvPassword))
			return
		}

		instance = newClient(values["username"], values["password"], values["token"], values["serverUrl"])
		instance.UsedLayer = usedLayer

		// Check if at least one authentication method is available
		if values["token"] == "" && (values["username"] == "" || values["password"] == "") {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("not logged in: use either environment variables %s or %s and %s, or use `ionosctl login`", constants.EnvToken, constants.EnvUsername, constants.EnvPassword))
			return
		}

		if err := instance.TestCreds(); err != nil {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("failed creating client: %w", err))
		}
	})

	return instance, getClientErr
}

// Must gets the client obj or fatally dies
// You can provide some optional custom error handlers as params. The err is sent to each error handler in order.
// The default error handler is die.Die which exits with code 1 and violently terminates the program
func Must(ehs ...func(error)) *Client {
	client, err := Get()
	if err != nil {
		if len(ehs) == 0 {
			die.Die(fmt.Errorf("failed getting client: %w", err).Error())
		} else {
			for _, eh := range ehs {
				eh(err)
			}
		}
	}
	return client
}

// NewClient bypasses the singleton check, not recommended for normal use.
func NewClient(name, pwd, token, hostUrl string) *Client {
	return newClient(name, pwd, token, hostUrl)
}

func TestCreds(user, pass, token string) error {
	cl := newClient(user, pass, token, constants.DefaultApiURL)
	return cl.TestCreds()
}

func (c *Client) TestCreds() error {
	_, _, err := c.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		usedScheme := "used token"
		if c.CloudClient.GetConfig().Token == "" {
			usedScheme = fmt.Sprintf("used username '%s' and password", c.CloudClient.GetConfig().Username)
		}
		return fmt.Errorf("credentials test failed. %s: %w", usedScheme, err)
	}

	return nil
}

// IsTokenAuth returns true if a token is being used for authentication. Otherwise, username & password were used.
func (c *Client) IsTokenAuth() bool {
	return c.CloudClient.GetConfig().Token != ""
}
