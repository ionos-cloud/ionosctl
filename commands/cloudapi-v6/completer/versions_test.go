package completer

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestK8sClusterUpgradeVersions(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	assert.NoError(t, os.Setenv(ionoscloud.IonosUsernameEnvVar, "user"))
	assert.NoError(t, os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass"))
	assert.NoError(t, os.Setenv(ionoscloud.IonosTokenEnvVar, "tok"))
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)

	buf := new(bytes.Buffer)
	K8sClusterUpgradeVersions(buf, "123")
	assert.True(t, regexp.MustCompile(`401 Unauthorized`).Match(buf.Bytes()))
}

func TestK8sNodePoolUpgradeVersions(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	assert.NoError(t, os.Setenv(ionoscloud.IonosUsernameEnvVar, "user"))
	assert.NoError(t, os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass"))
	assert.NoError(t, os.Setenv(ionoscloud.IonosTokenEnvVar, "tok"))
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)

	buf := new(bytes.Buffer)
	K8sNodePoolUpgradeVersions(buf, "123", "456")
	assert.True(t, regexp.MustCompile(`401 Unauthorized`).Match(buf.Bytes()))
}
