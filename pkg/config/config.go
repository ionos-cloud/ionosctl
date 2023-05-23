package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	cloudv6 "github.com/ionos-cloud/sdk-go/v6"

	"github.com/spf13/viper"
)

var FieldsWithSensitiveDataInConfigFile = []string{
	"userdata.name", "userdata.password", "userdata.token", // credentials stored in config file pre June 2023
	constants.Token, // credentials currently stored in config file
}

// GetUserData is deprecated
// It is hard to tell what values it will use and hence is un go-ish
// Use config.WriteFile
func GetUserData() map[string]string {
	return map[string]string{
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

func checkFilePermissions(fileInfo os.FileInfo, path string) error {
	var requiredPerm os.FileMode
	if runtime.GOOS == "windows" {
		requiredPerm = 0666
	} else {
		requiredPerm = 0600
	}

	if fileInfo.Mode().Perm() != requiredPerm {
		return fmt.Errorf("invalid permissions for %s: expected %o, got %o", path, requiredPerm, fileInfo.Mode().Perm())
	}
	return nil
}

// ReadFile reads config file at getConfigPath() and returns its data as a map
func ReadFile() (map[string]string, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed getting config file info: %w", err)
	}

	err = checkFilePermissions(fileInfo, path)
	if err != nil {
		return nil, fmt.Errorf("failed config file permissions check: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading config file: %w", err)
	}

	var result map[string]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling config file data: %w", err)
	}

	return result, nil
}

func LoadData(data map[string]string) {
	for key, value := range data {
		viper.Set(key, value)
	}
}

// LoadFile is a QOL func which does LoadData(ReadFile()) with err checking
// Exists because of old code
func LoadFile() error {
	data, err := ReadFile()
	if err != nil {
		return fmt.Errorf("failed reading config file: %w", err)
	}
	LoadData(data)
	return nil
}

// getConfigPath retrieves the configuration file path and makes it absolute if it isn't.
func getConfigPath() (string, error) {
	path := GetConfigFile()

	if !filepath.IsAbs(path) {
		// TODO: What is the point of turning this into an abs path ?
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	return path, nil
}

// Load binds environment variables (IONOS_USERNAME, IONOS_PASSWORD) to viper, and attempts
// to read config file for setting fallbacks for these newly-bound viper vars
func Load() (err error) {
	_ = viper.BindEnv(constants.Token, cloudv6.IonosTokenEnvVar)
	_ = viper.BindEnv(constants.ServerUrl, cloudv6.IonosApiUrlEnvVar)

	err = LoadFile() // Use config file as a fallback for any of the above variables. Could be used only for api-url

	if viper.IsSet(constants.Token) || (viper.IsSet("IONOS_USERNAME") && viper.IsSet("IONOS_PASSWORD")) {
		// Error thrown by LoadFile is recoverable in this case.
		// We don't want to throw an error e.g. if the user only uses the config file for api-url,
		// or if he has IONOS_TOKEN, or IONOS_USERNAME and IONOS_PASSWORD exported as env vars and no config file at all
		return nil
	}

	return fmt.Errorf("%w: Please export %s, or %s and %s, or do ionosctl login to generate a config file containing a token",
		err, cloudv6.IonosTokenEnvVar, cloudv6.IonosUsernameEnvVar, cloudv6.IonosPasswordEnvVar)
}

func WriteFile(data map[string]string) error {
	f, err := configFileWriter()
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := json.MarshalIndent(data, "", "  ")
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
