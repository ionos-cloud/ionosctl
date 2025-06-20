package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"

	"github.com/spf13/viper"
)

// GetServerUrl returns the server URL the SDK should use, with support for layered fallbacks.
func GetServerUrl() string {
	viper.AutomaticEnv()
	if val := viper.GetString(constants.ArgServerUrl); viper.IsSet(constants.ArgServerUrl) && val != "" {
		// 1. Above all, use global flag val
		return val
	}
	if val := viper.GetString(constants.EnvServerUrl); viper.IsSet(constants.EnvServerUrl) && val != "" {
		// 2. Fallback to non-empty env vars
		return val
	}

	return ""
}

// GetServerUrlOrApiIonos calls GetServerUrl and returns https://api.ionos.com if empty
//
// It is a useful func for informing the user of the behaviour of the SDKs - For the SDKs if the server URL is empty, they will default to https://api.ionos.com
func GetServerUrlOrApiIonos() string {
	if val := GetServerUrl(); val != "" {
		return val
	}
	return constants.DefaultApiURL
}

// GetConfigFilePath sanitizes the --config flag input and returns the path to the config file.
// If none set, it returns the default config path.
func GetConfigFilePath() string {
	path := filepath.Join(getConfigHomeDir(), constants.DefaultConfigFileName)
	if fn := constants.ArgConfig; viper.IsSet(fn) {
		path = viper.GetString(fn)
	}

	// We don't perform an `isAbs` check before turning it into an absolute path
	// because it internally has this check and will perform filepath.Clean on it if so
	// which is a great thing to have (sanitizes the path for multiple separators, path name elements, etc.)
	absPath, err := filepath.Abs(path)
	if err != nil {
		// just use the given provided by the user if err. Read and Write can still handle relative paths,
		// the only downside is annoyance for the user of not having his pwd prepended to `ionosctl location`
		return path
	}

	// Always prefer returning an absolute, cleaned path if possible.
	return absPath
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

// Read reads the config file and returns its data as a map
func Read() (map[string]string, error) {
	path := GetConfigFilePath()

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
