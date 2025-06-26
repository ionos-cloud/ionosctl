package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cfg "github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var once sync.Once
var instance *Client

func retrieveConfigFile() (*fileconfiguration.FileConfig, string, error) {
	// 1) --config flag
	if cfg, path, err := loadFromFlag(); cfg != nil || err != nil {
		return cfg, path, err
	}

	// 2) SDK env var
	if cfg, path, err := loadFromEnvVar(); cfg != nil || err != nil {
		return cfg, path, err
	}

	// 3) migrate old JSON if applicable
	if cfg, path, err := loadFromJSONMigration(); cfg != nil || err != nil {
		return cfg, path, err
	}

	// 4) default SDK path
	if cfg, path, err := loadFromSDKDefault(); cfg != nil || err != nil {
		return cfg, path, err
	}

	// note: if we reach this point, no config file was found
	// though old CLI behaviour was to return the default config path
	return nil, viper.GetString(constants.ArgConfig), nil
}

func Get() (*Client, error) {
	var getClientErr error

	// Every time, pick up the desired host URL from Viper
	desiredURL := viper.GetString(constants.ArgServerUrl)

	once.Do(
		func() {
			config, path, err := retrieveConfigFile()
			if err != nil {
				getClientErr = fmt.Errorf("failed to retrieve config file: %w", err)
				return
			}

			if instance == nil && os.Getenv(constants.EnvToken) != "" {
				instance = newClient("", "", os.Getenv(constants.EnvToken), desiredURL)
				instance.AuthSource = AuthSourceEnvBearer
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

func loadFromFlag() (*fileconfiguration.FileConfig, string, error) {
	path := viper.GetString(constants.ArgConfig)
	if path == "" {
		return nil, "", nil
	}
	cfg, err := fileconfiguration.New(path)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, path, fmt.Errorf("failed to load config from --config '%s': %w", path, err)
	}
	return cfg, path, nil
}

func loadFromEnvVar() (*fileconfiguration.FileConfig, string, error) {
	path := os.Getenv(shared.IonosFilePathEnvVar)
	if path == "" {
		return nil, "", nil
	}
	cfg, err := fileconfiguration.New(path)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, path, fmt.Errorf(
			"failed to load config from env var %s='%s': %w",
			shared.IonosFilePathEnvVar, path, err,
		)
	}
	return cfg, path, nil
}

func loadFromJSONMigration() (*fileconfiguration.FileConfig, string, error) {
	yamlPath := viper.GetString(constants.ArgConfig)
	if yamlPath == "" {
		return nil, "", nil
	}

	jsonPath := filepath.Join(filepath.Dir(yamlPath), "config.json")
	if _, err := os.Stat(jsonPath); err != nil {
		return nil, "", nil // no JSON to migrate
	}

	migrated, err := cfg.MigrateFromJSON(jsonPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed migrating %s → YAML: %w", jsonPath, err)
	}
	if migrated == nil {
		return nil, "", nil
	}

	out, _ := yaml.Marshal(migrated)
	if err := os.WriteFile(yamlPath, out, 0o600); err != nil {
		fmt.Fprintf(os.Stderr,
			"Warning: could not write migrated config to %s: %v\n",
			yamlPath, err,
		)
	}
	return migrated, yamlPath, nil
}

func loadFromSDKDefault() (*fileconfiguration.FileConfig, string, error) {
	defaultPath, err := fileconfiguration.DefaultConfigFileName()
	if err != nil {
		return nil, defaultPath, fmt.Errorf("failed to get default config path: %w", err)
	}
	cfg, err := fileconfiguration.New(defaultPath)
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		return nil, defaultPath, fmt.Errorf("failed to load default config '%s': %w", defaultPath, err)
	}
	return cfg, defaultPath, nil
}
