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

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	cloudv6 "github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/spf13/viper"
)

func GetUserData() map[string]string {
	return map[string]string{
		constants.Username:  viper.GetString(constants.Username),
		constants.Password:  viper.GetString(constants.Password),
		constants.Token:     viper.GetString(constants.Token),
		constants.ServerUrl: viper.GetString(constants.ServerUrl),
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
	if url := viper.GetString(constants.ServerUrl); url != "" {
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
		return statErr
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
	return viper.ReadInConfig()
}

// Load binds environment variables (IONOS_USERNAME, IONOS_PASSWORD) to viper, and attempts
// to read config file for setting fallbacks for these newly-bound viper vars
func Load() (err error) {
	_ = viper.BindEnv(constants.Username, cloudv6.IonosUsernameEnvVar)
	_ = viper.BindEnv(constants.Password, cloudv6.IonosPasswordEnvVar)
	_ = viper.BindEnv(constants.Token, cloudv6.IonosTokenEnvVar)
	_ = viper.BindEnv(constants.ServerUrl, cloudv6.IonosApiUrlEnvVar)

	err = LoadFile() // Use config file as a fallback for any of the above variables. Could be used only for api-url

	if viper.IsSet(constants.Token) || (viper.IsSet(constants.Username) && viper.IsSet(constants.Password)) {
		// Error thrown by LoadFile is recoverable in this case.
		// We don't want to throw an error e.g. if the user only uses the config file for api-url,
		// or if he has IONOS_TOKEN, or IONOS_USERNAME and IONOS_PASSWORD exported as env vars and no config file at all
		return nil
	}

	return fmt.Errorf("%w: Please export %s, or %s and %s, or do ionosctl login to generate a config file",
		err, cloudv6.IonosTokenEnvVar, cloudv6.IonosUsernameEnvVar, cloudv6.IonosPasswordEnvVar)
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
