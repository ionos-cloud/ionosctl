package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetUserData() map[string]string {
	return map[string]string{
		Username:     viper.GetString(Username),
		Password:     viper.GetString(Password),
		ArgServerUrl: viper.GetString(ArgServerUrl),
	}
}

func GetConfigFilePath() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return err.Error()
	}
	return filepath.Join(configPath, DefaultConfigFileName)
}

func LoadFile() error {
	viper.SetConfigFile(viper.GetString(ArgConfig))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
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
