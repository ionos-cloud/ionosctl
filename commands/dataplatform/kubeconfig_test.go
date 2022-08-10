package dataplatform

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testKubeconfigVar = "test-kubeconfig"
	testKubeconfigErr = errors.New("kubeconfig test error")
)

func TestRunKubeConfigGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testKubeconfigVar)
		rm.DataPlatformMocks.Cluster.EXPECT().GetKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, nil)
		err := RunKubeConfigGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunKubeConfigGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testKubeconfigVar)
		rm.DataPlatformMocks.Cluster.EXPECT().GetKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, testKubeconfigErr)
		err := RunKubeConfigGet(cfg)
		assert.Error(t, err)
	})
}
