package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testKubeconfigVar)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testKubeconfigVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(testKubeconfigVar, nil, testKubeconfigErr)
		err := RunK8sKubeconfigGet(cfg)
		assert.Error(t, err)
	})
}
