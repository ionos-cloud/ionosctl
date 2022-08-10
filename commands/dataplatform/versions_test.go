package dataplatform

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/stretchr/testify/assert"
)

var (
	testVersions    []string
	testVersionsErr = errors.New("test versions error")
)

func TestRunVersionsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultVersionsCols)
		rm.DataPlatformMocks.Versions.EXPECT().List().Return(testVersions, nil, nil)
		err := RunVersionsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVersionsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultVersionsCols)
		rm.DataPlatformMocks.Versions.EXPECT().List().Return(testVersions, nil, testVersionsErr)
		err := RunVersionsList(cfg)
		assert.Error(t, err)
	})
}
