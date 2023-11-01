package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testKubeconfigVar = "test-kubeconfig"
	testKubeconfigErr = errors.New("kubeconfig test error")
)

func TestRunK8sKubeconfigGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testKubeconfigVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, &testResponse, nil)
		err := RunK8sKubeconfigGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sKubeconfigGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testKubeconfigVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, testKubeconfigErr)
		err := RunK8sKubeconfigGet(cfg)
		assert.Error(t, err)
	})
}
