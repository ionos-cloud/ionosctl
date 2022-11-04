package config

import (
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetServerUrl(t *testing.T) {
	os.Clearenv()
	viper.Reset()

	// use env
	assert.NoError(t, os.Setenv(sdk.IonosApiUrlEnvVar, "url"))
	_ = Load() // ignore error since we just want to load the URL
	assert.Equal(t, "url", GetServerUrl())

	// from config
	os.Clearenv()
	viper.Reset()
	viper.SetConfigFile(filepath.Join("..", "testdata", "config.json"))
	viper.Set(constants.ArgConfig, filepath.Join("..", "testdata", "config.json"))
	assert.NoError(t, os.Chmod(filepath.Join("..", "testdata", "config.json"), 0600))
	assert.NoError(t, Load())
	assert.Equal(t, "https://api.ionos.com", GetServerUrl())

	viper.Reset()
	fs := pflag.NewFlagSet(constants.ArgServerUrl, pflag.ContinueOnError)
	_ = fs.String(constants.ArgServerUrl, "default", "test flag")
	viper.BindPFlags(fs)
	assert.Equal(t, "default", GetServerUrl())

	assert.NoError(t, fs.Parse([]string{"--" + constants.ArgServerUrl, "explicit"}))
	assert.Equal(t, "explicit", GetServerUrl())
}

func TestLoadFile(t *testing.T) {
	os.Clearenv()
	viper.Reset()

	viper.SetConfigFile(filepath.Join("..", "testdata", "config.json"))
	viper.Set(constants.ArgConfig, filepath.Join("..", "testdata", "config.json"))
	assert.NoError(t, os.Chmod(filepath.Join("..", "testdata", "config.json"), 0600))
	assert.NoError(t, LoadFile())
	assert.Equal(t, "test@ionos.com", viper.GetString(Username))
	assert.Equal(t, "test", viper.GetString(Password))
	assert.Equal(t, "jwt-token", viper.GetString(Token))
	assert.Equal(t, "https://api.ionos.com", viper.GetString(ServerUrl))
}

func TestLoadEnvFallback(t *testing.T) {
	os.Clearenv()
	viper.Reset()

	viper.SetConfigFile(filepath.Join("..", "testdata", "config.json"))
	viper.Set(constants.ArgConfig, filepath.Join("..", "testdata", "config.json"))
	assert.NoError(t, os.Setenv(sdk.IonosUsernameEnvVar, "user"))
	assert.NoError(t, os.Setenv(sdk.IonosPasswordEnvVar, "pass"))
	assert.NoError(t, os.Setenv(sdk.IonosTokenEnvVar, "token"))
	assert.NoError(t, os.Setenv(sdk.IonosApiUrlEnvVar, "url"))
	assert.NoError(t, Load())
	assert.Equal(t, "user", viper.GetString(Username))
	assert.Equal(t, "pass", viper.GetString(Password))
	assert.Equal(t, "token", viper.GetString(Token))
	assert.Equal(t, "url", viper.GetString(ServerUrl))
}

func TestRegressionLoadEnv(t *testing.T) {
	os.Clearenv()
	viper.Reset()

	os.Setenv(sdk.IonosUsernameEnvVar, "user")
	os.Setenv(sdk.IonosPasswordEnvVar, "pass")

	// This will fail if the config file if attempting to load the config file
	// Loading the config file is not necessary if environment variables are set
	assert.NoError(t, Load())

	assert.Equal(t, "user", viper.GetString(Username))
	assert.Equal(t, "pass", viper.GetString(Password))
}
