package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/spf13/viper"
)

var FieldsWithSensitiveDataInConfigFile = []string{
	constants.CfgUsername, constants.CfgPassword, constants.CfgToken,
}

func GetServerUrl() string {
	viper.AutomaticEnv()
	if flagVal := viper.GetString(constants.ArgServerUrl); flagVal != "" {
		// 1. Above all, use if the global flag is set
		return flagVal
	}
	if envVal := viper.GetString(constants.EnvServerUrl); envVal != "" {
		// 2. Fallback to non-empty env vars
		return envVal
	}
	if cfgVal := viper.GetString(constants.CfgServerUrl); cfgVal != "" {
		// 3. Fallback to non-empty cfg field
		return cfgVal
	}
	// 4. Return empty string. TODO: SDK (should) handle it. Test me
	return ""
	//return viper.GetString(constants.ArgServerUrl) // Return flag default. NOTE: DNS API uses a different set of URLs!
}

func GetConfigFile() string {
	if fn := constants.ArgConfig; viper.IsSet(fn) {
		return viper.GetString(fn)
	}
	return filepath.Join(getConfigHomeDir(), constants.DefaultConfigFileName)
}

func getConfigHomeDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		die.Die("is $HOME defined? couldn't get config dir: " + err.Error())
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

// Read reads config file at getConfigPath() and returns its data as a map
func Read() (map[string]string, error) {
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

func Write(data map[string]string) error {
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
	var filePath string
	if viper.IsSet(constants.ArgConfig) {
		filePath = viper.GetString(constants.ArgConfig)
	} else {
		configPath := getConfigHomeDir()
		err := os.MkdirAll(configPath, 0700) // Directory permissions are set to 0700
		if err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
		filePath = filepath.Join(configPath, constants.DefaultConfigFileName)
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600) // File is directly created with 0600 permissions
	if err != nil {
		return nil, fmt.Errorf("failed to create config file: %w", err)
	}

	return f, nil
}
