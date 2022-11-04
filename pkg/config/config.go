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

	"github.com/ionos-cloud/ionosctl/pkg/constants"

	sdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func GetUserData() map[string]string {
	return map[string]string{
		Username:  viper.GetString(Username),
		Password:  viper.GetString(Password),
		Token:     viper.GetString(Token),
		ServerUrl: viper.GetString(ServerUrl),
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
	if viper.IsSet(constants.ArgServerUrl) {
		return viper.GetString(constants.ArgServerUrl)
	}
	if url := viper.GetString(ServerUrl); url != "" {
		return url
	}
	return viper.GetString(constants.ArgServerUrl)
}

func GetConfigFile() string {
	return filepath.Join(getConfigHomeDir(), constants.DefaultConfigFileName)
}

func getConfigHomeDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return err.Error()
	}
	return filepath.Join(configPath, "ionosctl")
}

func getPermsByOS(os string) int {
	if os == "windows" {
		return 666
	} else {
		return 600
	}
}

func LoadFile() error {
	path := viper.GetString(constants.ArgConfig)
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return fmt.Errorf("error getting credentials: nor $%s, $%s, $%s set, nor config file: %s",
			sdk.IonosUsernameEnvVar, sdk.IonosPasswordEnvVar, sdk.IonosTokenEnvVar, statErr.Error())
	}

	perm := fileInfo.Mode().Perm()
	permNumberBase10 := int64(perm)
	strBase10 := strconv.FormatInt(permNumberBase10, 8)
	permNumber, _ := strconv.Atoi(strBase10)

	system := runtime.GOOS

	permNumberExpected := getPermsByOS(system)
	if permNumber != permNumberExpected {
		return fmt.Errorf("config file %s has wrong permissions: %d, should be %d", path, permNumber, permNumberExpected)
	}

	viper.SetConfigFile(viper.GetString(constants.ArgConfig))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// Load collects config data from the config file, using environment variables as fallback.
func Load() (err error) {
	_ = viper.BindEnv(Username, sdk.IonosUsernameEnvVar)
	_ = viper.BindEnv(Password, sdk.IonosPasswordEnvVar)
	_ = viper.BindEnv(Token, sdk.IonosTokenEnvVar)
	_ = viper.BindEnv(ServerUrl, sdk.IonosApiUrlEnvVar)

	if viper.GetString(Username) == "" && viper.GetString(Token) == "" {
		err = LoadFile()
	}

	return
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
	if !viper.IsSet(constants.ArgConfig) {
		configPath := getConfigHomeDir()
		err := os.MkdirAll(configPath, 0755)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.Create(viper.GetString(constants.ArgConfig))
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(viper.GetString(constants.ArgConfig), 0600); err != nil {
		return nil, err
	}
	return f, nil
}
