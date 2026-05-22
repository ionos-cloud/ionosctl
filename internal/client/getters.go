package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
)

var (
	once         sync.Once
	instance     *Client
	getClientErr error
)

func Get() (*Client, error) {
	once.Do(
		func() {
			src, err := retrieveConfigFile()
			if err != nil {
				getClientErr = fmt.Errorf("failed to retrieve config file: %w", err)
				return
			}
			config, path := src.Config, src.Path

			if instance == nil && os.Getenv(constants.EnvToken) != "" {
				instance = newClient("", "", os.Getenv(constants.EnvToken), config)
				instance.AuthSource = AuthSourceEnvBearer
			}

			if instance == nil && os.Getenv(constants.EnvUsername) != "" && os.Getenv(constants.EnvPassword) != "" {
				instance = newClient(os.Getenv(constants.EnvUsername), os.Getenv(constants.EnvPassword), "", config)
				instance.AuthSource = AuthSourceEnvBasic
			}

			if instance == nil && config != nil && config.GetCurrentProfile() != nil &&
				config.GetCurrentProfile().Credentials.Token != "" {
				instance = newClient("", "", config.GetCurrentProfile().Credentials.Token, config)
				instance.AuthSource = AuthSourceCfgBearer
			}

			if instance == nil && config != nil && config.GetCurrentProfile() != nil &&
				config.GetCurrentProfile().Credentials.Username != "" &&
				config.GetCurrentProfile().Credentials.Password != "" {
				instance = newClient(
					config.GetCurrentProfile().Credentials.Username,
					config.GetCurrentProfile().Credentials.Password,
					"", config)
				instance.AuthSource = AuthSourceCfgBasic
			}

			if instance == nil {
				instance = newClient("", "", "", config)
				instance.AuthSource = AuthSourceNone
				getClientErr = fmt.Errorf("no credentials found, please update your config file at "+
					"'ionosctl cfg location', or generate a new one with 'ionosctl login', "+
					"or set the environment variable %s or %s and %s",
					constants.EnvToken, constants.EnvUsername, constants.EnvPassword)
				return
			}

			instance.Config = config
			instance.ConfigPath = path
		},
	)

	// No URL-change rebuild. Each SDK client is created with its own correct URL.
	// Regional commands use NewRegionalConfig() for location-specific calls.

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

// NewClient creates a new client with the given credentials and no config file.
// It is used for testing purposes or when you want to create a client with specific credentials.
// It is highly recommended to use the Get() or Must() functions instead, as they handle the configuration file and environment variables automatically.
func NewClient(name, pwd, token string) *Client {
	return newClient(name, pwd, token, nil)
}

func TestCreds(user, pass, token string) error {
	cl := newClient(user, pass, token, nil)
	return cl.TestCreds()
}

func (c *Client) TestCreds() error {
	_, _, err := c.CloudClient.DefaultApi.ApiInfoGet(context.Background()).MaxResults(1).Depth(0).Execute()
	if err != nil {
		return fmt.Errorf("credentials test failed. used %s: %w", c.AuthSource, err)
	}

	return nil
}

// EnforceClient sets the global client instance to a new client with the given credentials (
// use only for testing/special cases)
func EnforceClient(user, pass, token string) {
	instance = newClient(user, pass, token, nil)
}
