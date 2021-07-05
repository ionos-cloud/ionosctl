package config

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ionos-cloud/sdk-go/v5"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/viper"
)

func GetUserData() map[string]string {
	return map[string]string{
		Username: viper.GetString(Username),
		Password: viper.GetString(Password),
		Token:    viper.GetString(Token),
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
	path := viper.GetString(ArgConfig)
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	system := runtime.GOOS
	fmt.Printf("SYSTEEEEEEEEEEMMMMMMMMM: %v", system)
	//fmt.Println("path: " + path)
	fileInfo, statErr := os.Stat(path)
	if statErr != nil {
		return statErr
	}

	perm := fileInfo.Mode().Perm()
	permNumberBase10 := int64(perm)
	strBase10 := strconv.FormatInt(permNumberBase10, 8)
	permNumber, _ := strconv.Atoi(strBase10)

	if permNumber == int(600) {
		viper.SetConfigFile(viper.GetString(ArgConfig))
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}
		return nil
	} else {
		fmt.Printf("perm: %v", permNumber)
		return errors.New("no permission for the config file, expected 600")
	}

}

// Load collects config data from the config file, using environment variables as fallback.
func Load() (err error) {
	_ = viper.BindEnv(Username, sdk.IonosUsernameEnvVar)
	_ = viper.BindEnv(Password, sdk.IonosPasswordEnvVar)
	_ = viper.BindEnv(Token, sdk.IonosTokenEnvVar)

	if viper.GetString(Username) == "" || viper.GetString(Password) == "" || viper.GetString(Token) == "" {
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
