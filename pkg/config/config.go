package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	sdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

type UserData struct {
	Username string `mapstructure:"name,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Token string `mapstructure:"token,omitempty"`
	ServerUrl string `mapstructure:"api-url"`
}

type Config struct {
	Data UserData `mapstructure:"userdata"`
}

func GetUserData() Config {
	return Config {
		Data: UserData {
			Username:  viper.GetString(Username),
			Password:  viper.GetString(Password),
			Token:     viper.GetString(Token),
			ServerUrl: viper.GetString(ServerUrl),
		},
	}
}

// GetServerUrl returns the API URL from flags, config or env in order of priority.
// The caller must ensure to load config or env vars beforehand, so they can be included.
//
// Priority:
// 1. Explicit flag
// 2. Env/config file
// 3. Flag default value
func GetServerUrl() string {
	if viper.IsSet(ArgServerUrl) {
		return viper.GetString(ArgServerUrl)
	}
	if url := viper.GetString(ServerUrl); url != "" {
		return url
	}
	return viper.GetString(ArgServerUrl)
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
	path := viper.GetString(ArgConfig)
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return errors.New(fmt.Sprintf("error getting credentials: nor $%s, $%s, $%s set, nor config file: %s",
			sdk.IonosUsernameEnvVar, sdk.IonosPasswordEnvVar, sdk.IonosTokenEnvVar, statErr.Error()))
	}

	perm := fileInfo.Mode().Perm()
	permNumberBase10 := int64(perm)
	strBase10 := strconv.FormatInt(permNumberBase10, 8)
	permNumber, _ := strconv.Atoi(strBase10)

	system := runtime.GOOS
	if system == "windows" {
		if permNumber == int(666) {
			viper.SetConfigFile(viper.GetString(ArgConfig))
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("no permission for the config file, expected 600")
		}
	} else {
		if permNumber == int(600) {
			viper.SetConfigFile(viper.GetString(ArgConfig))
			err := viper.ReadInConfig()
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("no permission for the config file, expected 600")
		}
	}

}

// Load collects config data from the config file, using environment variables as fallback.
func Load() (err error) {
	_ = viper.BindEnv(Username, sdk.IonosUsernameEnvVar)
	_ = viper.BindEnv(Password, sdk.IonosPasswordEnvVar)
	_ = viper.BindEnv(Token, sdk.IonosTokenEnvVar)
	_ = viper.BindEnv(ServerUrl, sdk.IonosApiUrlEnvVar)

	if viper.GetString(Username) == "" && viper.GetString(Token) == "" {
		if err = LoadFile(); err != nil {
			return err
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
