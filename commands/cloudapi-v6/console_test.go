package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testConsole = resources.RemoteConsoleUrl{
		RemoteConsoleUrl: ionoscloud.RemoteConsoleUrl{
			Url: &testConsoleVar,
		},
	}
	testConsoleVar = "test-console"
	testConsoleErr = errors.New("console test: error occurred")
)

func TestRunServerConsoleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testConsoleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testConsoleVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetRemoteConsoleUrl(testConsoleVar, testConsoleVar).Return(testConsole, nil, nil)
		err := RunServerConsoleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerConsoleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testConsoleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testConsoleVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetRemoteConsoleUrl(testConsoleVar, testConsoleVar).Return(testConsole, nil, testConsoleErr)
		err := RunServerConsoleGet(cfg)
		assert.Error(t, err)
	})
}
