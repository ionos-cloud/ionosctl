package configgen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/viper"
)

// utils for file operations

func Location() string {
	if viper.IsSet(constants.ArgConfig) {
		return viper.GetString(constants.ArgConfig)
	}
	configPath := getConfigHomeDir()
	return filepath.Join(configPath, constants.DefaultConfigFileName)
}

func configFileWriter() (io.WriteCloser, error) {
	var filePath string
	if viper.IsSet(constants.ArgConfig) {
		filePath = viper.GetString(constants.ArgConfig)
	} else {
		filePath = Location()
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory structure '%s': %w", filePath, err)
		}
	}

	err := checkConfigExistsAndAskForReplace(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed handling existing config file '%s': %w", filePath, err)
	}

	f, err := os.OpenFile(
		filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600,
	) // File is directly created with 0600 permissions
	if err != nil {
		return nil, fmt.Errorf("failed to create config file '%s': %w", filePath, err)
	}

	return f, nil
}

func getConfigHomeDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		configPath = os.Getenv("HOME")
	}
	return filepath.Join(configPath, "ionosctl")
}

func checkConfigExistsAndAskForReplace(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	if !confirm.FAsk(os.Stdin,
		fmt.Sprintf("Config file %s already exists. Do you want to replace it", configPath),
		viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("error deleting config file %s: %w", configPath, err)
	}

	return nil
}
