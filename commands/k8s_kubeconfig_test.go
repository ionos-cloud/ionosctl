package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testKubeconfigVar)
		rm.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, nil)
		err := RunK8sKubeconfigGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sKubeconfigGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testKubeconfigVar)
		rm.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, testKubeconfigErr)
		err := RunK8sKubeconfigGet(cfg)
		assert.Error(t, err)
	})
}
