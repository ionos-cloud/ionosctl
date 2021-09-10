package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testToken = resources.Token{
		Token: ionoscloud.Token{
			Token: &testTokenVar,
		},
	}
	testTokenVar = "test-token"
	testTokenErr = errors.New("token test: error occurred")
)

func TestRunServerTokenGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgServerId), testTokenVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetToken(testTokenVar, testTokenVar).Return(testToken, nil, nil)
		err := RunServerTokenGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerTokenGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgServerId), testTokenVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetToken(testTokenVar, testTokenVar).Return(testToken, nil, testTokenErr)
		err := RunServerTokenGet(cfg)
		assert.Error(t, err)
	})
}
