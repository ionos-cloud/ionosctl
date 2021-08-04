package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testToken = v6.Token{
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testTokenVar)
		rm.Server.EXPECT().GetToken(testTokenVar, testTokenVar).Return(testToken, nil, nil)
		err := RunServerTokenGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerTokenGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testTokenVar)
		rm.Server.EXPECT().GetToken(testTokenVar, testTokenVar).Return(testToken, nil, testTokenErr)
		err := RunServerTokenGet(cfg)
		assert.Error(t, err)
	})
}
