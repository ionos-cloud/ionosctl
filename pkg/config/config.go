package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
)

func WriteConfigFile() error {
	f, err := configFileWriter()
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return errors.New("unable to encode configuration to YAML format")
	}

	_, err = f.Write(b)
	if err != nil {
		return errors.New("unable to write configuration")
	}

	return nil
}

func configFileWriter() (io.WriteCloser, error) {
	cfgFile := viper.GetString(ArgConfig)

	f, err := os.Create(cfgFile)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(cfgFile, 0600); err != nil {
		return nil, err
	}

	return f, nil
}

func GetConfigFilePath() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return err.Error()
	}
	return filepath.Join(configPath, DefaultConfigFileName)
}

func GetAPIClient() (*ionoscloud.APIClient, error) {
	viper.SetConfigFile(viper.GetString(ArgConfig))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if viper.GetString(Username) == "" {
		return nil, errors.New("no username set")
	}
	if viper.GetString(Password) == "" {
		return nil, errors.New("no password set")
	}

	cfg := GetAPIClientConfig()
	return ionoscloud.NewAPIClient(cfg), nil
}

func GetAPIClientConfig() *ionoscloud.Configuration {
	return &ionoscloud.Configuration{
		Username: viper.GetString(Username),
		Password: viper.GetString(Password),
		Servers: ionoscloud.ServerConfigurations{
			ionoscloud.ServerConfiguration{
				URL: viper.GetString(ArgServerUrl),
			},
		},
	}
}
