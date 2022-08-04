package completer

import (
	"bufio"
	"bytes"
	"os"
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
	// ToDo: add this again when data platform in prod
	//re := regexp.MustCompile(`401 Unauthorized`)
	//assert.True(t, re.Match(b.Bytes()))
}
