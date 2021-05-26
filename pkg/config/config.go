package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func GetUserData() map[string]string {
	return map[string]string{
		Username:     viper.GetString(Username),
		Password:     viper.GetString(Password),
		Token:        viper.GetString(Token),
		ArgServerUrl: viper.GetString(ArgServerUrl),
	}
}

func GetConfigFile() string {
	return filepath.Join(getConfigHomeDir(), DefaultConfigFileName)
}

func getConfigHomeDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return err.Error()
	}
	return filepath.Join(configPath, "ionosctl")
}

func LoadFile() error {
	viper.SetConfigFile(viper.GetString(ArgConfig))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// Load collects config data from the config file, using environment variables as fallback.
func Load() (err error) {
	if err = LoadFile(); err != nil {
		pathErr := &os.PathError{}
		if errors.As(err, &viper.ConfigFileNotFoundError{}) || errors.As(err, &pathErr) {
			_ = viper.BindEnv(Username, ionoscloud.IonosUsernameEnvVar)
			_ = viper.BindEnv(Password, ionoscloud.IonosPasswordEnvVar)
			_ = viper.BindEnv(Token, ionoscloud.IonosTokenEnvVar)
			return nil
		}
	}
	return err
}

func WriteFile() error {
	f, err := configFileWriter()
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := json.MarshalIndent(GetUserData(), "", "  ")
	if err != nil {
		return errors.New("unable to encode configuration to JSON format")
	}

	_, err = f.Write(b)
	if err != nil {
		return errors.New("unable to write configuration")
	}
	return nil
}

func configFileWriter() (io.WriteCloser, error) {
	if !viper.IsSet(ArgConfig) {
		configPath := getConfigHomeDir()
		err := os.MkdirAll(configPath, 0755)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.Create(viper.GetString(ArgConfig))
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(viper.GetString(ArgConfig), 0600); err != nil {
		return nil, err
	}
	return f, nil
}
