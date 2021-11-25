package completer

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetTokensIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(sdkgoauth.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(sdkgoauth.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(sdkgoauth.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	TokensIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
