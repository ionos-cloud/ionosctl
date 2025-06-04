package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	cfg "github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var once sync.Once
var instance *Client

// retrieveConfigFile tries to retrieve the configuration file, in this order:
// 1. From the --config flag (default value or set by the user)
// 2. From the fileconfiguration.NewFromEnv result, if not nil
// 3. From the fileconfiguration.New("") default, if not nil
// NOTE: If the config file is not found in either location, it will simply return nil without an error.
func retrieveConfigFile() (*fileconfiguration.FileConfig, error) {
	var config *fileconfiguration.FileConfig

	// --- try the --config flag first
	if viper.GetString(constants.ArgConfig) != "" {
		config, err := fileconfiguration.New(viper.GetString(constants.ArgConfig))
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			// only return an error if the config file exists but is invalid
			return nil, fmt.Errorf("failed to create config from --config flag: %w", err)
		}
		if config != nil {
			return config, nil
		}
	}

	// --- try the config file from the sdk env var

	config, err := fileconfiguration.NewFromEnv()
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		// only return an error if the config file exists but is invalid
		return nil, fmt.Errorf("failed to create config from env: %w", err)
	}
	if config != nil {
		return config, nil
	}

	// --- try the default sdk path

	defaultSdkConfigPath, err := fileconfiguration.DefaultConfigFileName()
	if err != nil {
		return nil, fmt.Errorf("failed to get default config file path: %w", err)
	}
	config, err = fileconfiguration.New(defaultSdkConfigPath)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		// only return an error if the config file exists but is invalid
		return nil, fmt.Errorf("failed to create default config from %s: %w",
			defaultSdkConfigPath, err)
	}

	return config, nil
}

func Get() (*Client, error) {
	var getClientErr error

	// Every time, pick up the desired host URL from Viper
	desiredURL := viper.GetString(constants.ArgServerUrl)

	once.Do(
		func() {
			config, err := retrieveConfigFile()
			if err != nil {
				getClientErr = fmt.Errorf("failed to retrieve config file: %w", err)
				return
			}

			if os.Getenv(constants.EnvToken) != "" {
				instance = newClient("", "", os.Getenv(constants.EnvToken), "")
				instance.AuthSource = AuthSourceEnvBearer
				return
			}

			if os.Getenv(constants.EnvUsername) != "" && os.Getenv(constants.EnvPassword) != "" {
				instance = newClient(os.Getenv(constants.EnvUsername), os.Getenv(constants.EnvPassword), "", "")
				instance.AuthSource = AuthSourceEnvBasic
			}

			instance.Config = config
			instance.AuthSource = AuthSourceCfgBearer
			if instance.Config == nil {
				getClientErr = fmt.Errorf("no configuration file found, please use 'ionosctl login' "+
					"or set the environment variable %s or %s and %s",
					constants.EnvToken, constants.EnvUsername, constants.EnvPassword)
				return
			}

			if instance == nil && os.Getenv(constants.EnvUsername) != "" && os.Getenv(constants.EnvPassword) != "" {
				instance = newClient(os.Getenv(constants.EnvUsername), os.Getenv(constants.EnvPassword), "", desiredURL)
				instance.AuthSource = AuthSourceEnvBasic
			}

			if instance == nil && config.GetCurrentProfile() != nil &&
				config.GetCurrentProfile().Credentials.Token != "" {
				instance = newClient("", "", config.GetCurrentProfile().Credentials.Token, desiredURL)
				instance.AuthSource = AuthSourceCfgBearer
			}

			if instance == nil && config.GetCurrentProfile() != nil &&
				config.GetCurrentProfile().Credentials.Username != "" &&
				config.GetCurrentProfile().Credentials.Password != "" {
				instance = newClient(
					config.GetCurrentProfile().Credentials.Username,
					config.GetCurrentProfile().Credentials.Password,
					"", desiredURL)
				instance.AuthSource = AuthSourceCfgBasic
			}

			if instance == nil {
				instance = newClient("", "", "", desiredURL)
				instance.AuthSource = AuthSourceNone
				getClientErr = fmt.Errorf("no credentials found, please update your config file at "+
					"'ionosctl cfg location', or generate a new one with 'ionosctl login', "+
					"or set the environment variable %s or %s and %s",
					constants.EnvToken, constants.EnvUsername, constants.EnvPassword)
			}

			instance.Config = config
			instance.ConfigPath = path
		},
	)

	// If we already have an instance, but the desiredURL has changed, rebuild it
	if instance != nil && instance.URLOverride != desiredURL {
		// preserve credentials / auth source
		name, pwd, token := instance.CloudClient.GetConfig().Username, instance.CloudClient.GetConfig().Password, instance.CloudClient.GetConfig().Token
		newInst := newClient(name, pwd, token, desiredURL)
		newInst.AuthSource = instance.AuthSource
		newInst.Config = instance.Config
		instance = newInst
		instance.URLOverride = desiredURL
	}

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

// NewClient creates a new client with the given credentials.
// It is used for testing purposes or when you want to create a client with specific credentials.
// It is highly recommended to use the Get() or Must() functions instead, as they handle the configuration file and environment variables automatically.
func NewClient(name, pwd, token, hostUrl string) *Client {
	return newClient(name, pwd, token, hostUrl)
}

func TestCreds(user, pass, token string) error {
	cl := newClient(user, pass, token, constants.DefaultApiURL)
	return cl.TestCreds()
}

func (c *Client) TestCreds() error {
	if c.URLOverride != constants.DefaultApiURL && c.URLOverride != "" {
		return nil
	}
	_, _, err := c.CloudClient.DefaultApi.ApiInfoGet(context.Background()).MaxResults(1).Depth(0).Execute()
	if err != nil {
		return fmt.Errorf("credentials test failed. used %s: %w", c.AuthSource, err)
	}

	return nil
}

// EnforceClient sets the global client instance to a new client with the given credentials (
// use only for testing/special cases)
func EnforceClient(user, pass, token, hostUrl string) {
	instance = newClient(user, pass, token, hostUrl)
}
