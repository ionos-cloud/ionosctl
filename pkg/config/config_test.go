package config

import (
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	viper.Reset()
	viper.SetConfigFile(filepath.Join("..", "testdata", "config.json"))
	viper.Set(ArgConfig, filepath.Join("..", "testdata", "config.json"))
	err := os.Chmod(filepath.Join("..", "testdata", "config.json"), 0600)
	assert.NoError(t, err)
	assert.NoError(t, LoadFile())
	assert.Equal(t, "test@ionos.com", viper.GetString(Username))
	assert.Equal(t, "test", viper.GetString(Password))
	assert.Equal(t, "jwt-token", viper.GetString(Token))
}

func TestLoadEnvFallback(t *testing.T) {
	viper.Reset()
	err := os.Setenv(sdk.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdk.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdk.IonosTokenEnvVar, "token")
	assert.NoError(t, err)
	assert.NoError(t, Load())
	assert.Equal(t, "user", viper.GetString(Username))
	assert.Equal(t, "pass", viper.GetString(Password))
	assert.Equal(t, "token", viper.GetString(Token))

	viper.Reset()
	viper.SetConfigFile("notfound.json")
	assert.NoError(t, Load())
	assert.Equal(t, "user", viper.GetString(Username))
	assert.Equal(t, "pass", viper.GetString(Password))
	assert.Equal(t, "token", viper.GetString(Token))
}
