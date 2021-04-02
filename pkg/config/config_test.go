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
	assert.NoError(t, Load())
	assert.Equal(t, "test@ionos.com", viper.GetString(Username))
	assert.Equal(t, "test", viper.GetString(Password))
	assert.Equal(t, "jwt-token", viper.GetString(Token))
}

func TestLoadEnvFallback(t *testing.T) {
	viper.Reset()
	os.Setenv(sdk.IonosUsernameEnvVar, "user")
	os.Setenv(sdk.IonosPasswordEnvVar, "pass")
	os.Setenv(sdk.IonosTokenEnvVar, "token")
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
