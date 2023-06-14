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
	credsProvenance *AuthProvenance

	CloudClient        *cloudv6.APIClient
	AuthClient         *sdkgoauth.APIClient
	CertManagerClient  *certmanager.APIClient
	PostgresClient     *postgres.APIClient
	MongoClient        *mongo.APIClient
	DataplatformClient *dataplatform.APIClient
	RegistryClient     *registry.APIClient
}

type AuthProvenance struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
	ApiUrl   string `json:"api_url"`
}

func (c *Client) GetProvenance() *AuthProvenance {
	if c.credsProvenance != nil {
		return c.credsProvenance
	}

	// Get / Must were skipped via NewClient, which happens for testing pkgs, or for testing non-active credentials, etc.
	return &AuthProvenance{
		Token:    "direct",
		Username: "direct",
		Password: "direct",
		ApiUrl:   "direct",
	}

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

		prov := AuthProvenance{}
		// Credentials and API URL priority: command line arguments -> environment variables -> config file
		token := viper.GetString(constants.ArgToken)
		prov.Token = "flag"
		if token == "" {
			token = viper.GetString(constants.EnvToken)
			prov.Token = "env"
		}
		if token == "" {
			token = viper.GetString(constants.CfgToken)
			prov.Token = "cfg"
		}

		hostUrl := viper.GetString(constants.ArgServerUrl)
		prov.ApiUrl = "flag"
		if hostUrl == "" {
			hostUrl = viper.GetString(constants.EnvServerUrl)
			prov.ApiUrl = "env"
		}
		if hostUrl == "" {
			hostUrl = viper.GetString(constants.CfgServerUrl)
			prov.ApiUrl = "cfg"
		}

		username := viper.GetString(constants.EnvUsername)
		prov.Username = "env"
		if username == "" {
			// Since June 2023 these config variables are no longer stored on successful `ionosctl login`.
			// However, we continue supporting them to avoid breaking changes. The user can manually add them to their config.
			username = viper.GetString(constants.CfgUsername)
			prov.Username = "cfg"
		}

		password := viper.GetString(constants.EnvPassword)
		prov.Password = "env"
		if password == "" {
			// Since June 2023 these config variables are no longer stored on successful `ionosctl login`.
			// However, we continue supporting them to avoid breaking changes. The user can manually add them to their config.
			password = viper.GetString(constants.CfgPassword)
			prov.Password = "cfg"
		}

		// Check if at least one authentication method is available
		if token == "" && (username == "" || password == "") {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("not logged in: use either environment variables %s or %s and %s, or use `ionosctl login`", constants.EnvToken, constants.EnvUsername, constants.EnvPassword))
			return
		}

		instance, err = newClient(username, password, token, hostUrl)
		instance.credsProvenance = &prov
		if err != nil {
			getClientErr = errors.Join(getClientErr, fmt.Errorf("failed creating client: %w", err))
		}

		err = TestCreds(username, password, token)
		if err != nil {
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
