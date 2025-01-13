package client

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
	"github.com/spf13/viper"
)

var once sync.Once
var instance *Client

func selectAuthLayer(layers []Layer) (values map[string]string, usedLayer Layer, err error) {
	for _, layer := range layers {
		token := viper.GetString(layer.TokenKey)
		username := viper.GetString(layer.UsernameKey)
		password := viper.GetString(layer.PasswordKey)

		if token != "" || (username != "" && password != "") {
			return map[string]string{
				"token":     token,
				"username":  username,
				"password":  password,
				"serverUrl": config.GetServerUrl(),
			}, layer, nil
		}
	}
	return nil, Layer{}, fmt.Errorf("none of the layers provided a value for either token or username & password. use `ionosctl whoami --provenance` for help")
}

// Get a client and possibly fail. Uses viper to get the credentials and API URL.
// The returned client is guaranteed to have working credentials
// Order:
// Explicit flags ( e.g. --token )
// Environment Variables ( e.g. IONOS_TOKEN )
// Config File ( e.g. userdata.token )
func Get() (*Client, error) {
	var getClientErr error

	once.Do(
		func() {
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

			values, usedLayer, err := selectAuthLayer(ConfigurationPriorityRules)
			if err != nil {
				getClientErr = errors.Join(getClientErr, fmt.Errorf("failed selecting an auth layer: %w", err))
				return
			}

			instance = newClient(
				values["username"], values["password"], values["token"], values["serverUrl"], &usedLayer,
			)

			if err := instance.TestCreds(); err != nil {
				getClientErr = errors.Join(getClientErr, fmt.Errorf("failed creating client: %w", err))
			}
		},
	)

	return instance, getClientErr
}

var MustDefaultErrHandler = func(err error) {
	die.Die(fmt.Errorf("failed getting client: %w", err).Error())
}

// Must gets the client obj or fatally dies
// You can provide some optional custom error handlers as params. The err is sent to each error handler in order.
// The default error handler is die.Die which exits with code 1 and violently terminates the program
func Must(ehs ...func(error)) *Client {
	client, err := Get()
	if err != nil {
		if len(ehs) > 0 {
			// Developer set custom err handlers (e.g. don't die, simply warn, etc)
			for _, eh := range ehs {
				eh(err)
			}
		} else {
			// Default error handler if none set
			MustDefaultErrHandler(err)
		}
	}
	return client
}

// NewClient bypasses the singleton check, not recommended for normal use.
func NewClient(name, pwd, token, hostUrl string) *Client {
	return newClient(name, pwd, token, hostUrl, nil)
}

func NewClientFromCfgData(values map[string]string) *Client {
	return newClient(
		values[constants.CfgUsername],
		values[constants.CfgPassword],
		values[constants.CfgToken],
		values[constants.CfgServerUrl],
		nil,
	)
}

func TestCreds(user, pass, token string) error {
	cl := newClient(user, pass, token, constants.DefaultApiURL, nil)
	return cl.TestCreds()
}

func (c *Client) TestCreds() error {
	if config.GetServerUrlOrApiIonos() != constants.DefaultApiURL {
		// TODO: Remove this if we can somehow reliably test credentials on all APIs.
		// TODO: This currently skips if the server URL is manually overwritten. (i.e. staging environment, or regional APIs)
		return nil
	}
	_, _, err := c.CloudClient.DefaultApi.ApiInfoGet(context.Background()).MaxResults(1).Depth(0).Execute()
	if err != nil {
		usedScheme := "used token"
		if c.CloudClient.GetConfig().Token == "" {
			usedScheme = fmt.Sprintf("used username '%s' and password", c.CloudClient.GetConfig().Username)
		}
		return fmt.Errorf("credentials test failed. %s: %w", usedScheme, err)
	}

	return nil
}

// EnforceClient sets the global client instance to a new client with the given credentials (
// use only for testing/special cases)
func EnforceClient(user, pass, token, hostUrl string) {
	instance = newClient(user, pass, token, hostUrl, nil)
}
