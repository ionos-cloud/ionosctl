package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	kubeconfigTestGet = resources.K8sKubeconfig{
		KubernetesConfig: ionoscloud.KubernetesConfig{
			Id: &testKubeconfigVar,
			Properties: &ionoscloud.KubernetesConfigProperties{
				Kubeconfig: &testKubeconfigVar,
			},
		},
	}
	testKubeconfigVar = "test-kubeconfig"
	testKubeconfigErr = errors.New("kubeconfig test error")
)

func TestRunK8sKubeconfigGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testKubeconfigVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(kubeconfigTestGet, &testResponse, nil)
		err := RunK8sKubeconfigGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sKubeconfigGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testKubeconfigVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(kubeconfigTestGet, nil, testKubeconfigErr)
		err := RunK8sKubeconfigGet(cfg)
		assert.Error(t, err)
	})
}
