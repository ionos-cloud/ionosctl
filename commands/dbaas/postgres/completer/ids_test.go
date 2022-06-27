package completer

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testIdVar = "test-id"

func TestClustersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(sdkgo.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosTokenEnvVar, "token")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	ClustersIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestBackupsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(sdkgo.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosTokenEnvVar, "token")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	BackupsIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestBackupsIdsForCluster(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(sdkgo.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosTokenEnvVar, "token")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	BackupsIdsForCluster(w, testIdVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPostgresVersions(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(sdkgo.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdkgo.IonosTokenEnvVar, "token")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	PostgresVersions(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
