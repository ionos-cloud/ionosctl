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
	testConsole = v6.RemoteConsoleUrl{
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testConsoleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testConsoleVar)
		rm.Server.EXPECT().GetRemoteConsoleUrl(testConsoleVar, testConsoleVar).Return(testConsole, nil, nil)
		err := RunServerConsoleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerConsoleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testConsoleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testConsoleVar)
		rm.Server.EXPECT().GetRemoteConsoleUrl(testConsoleVar, testConsoleVar).Return(testConsole, nil, testConsoleErr)
		err := RunServerConsoleGet(cfg)
		assert.Error(t, err)
	})
}
