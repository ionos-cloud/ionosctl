package commands

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	kubeconfigTest = resources.K8sKubeconfig{
		KubernetesConfig: ionoscloud.KubernetesConfig{
			Properties: &ionoscloud.KubernetesConfigProperties{
				Kubeconfig: &testKubeconfigVar,
			},
		},
	}
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testKubeconfigVar)
		rm.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(kubeconfigTestGet, nil, nil)
		err := RunK8sKubeconfigGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sKubeconfigGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testKubeconfigVar)
		rm.K8s.EXPECT().ReadKubeConfig(testKubeconfigVar).Return(kubeconfigTestGet, nil, testKubeconfigErr)
		err := RunK8sKubeconfigGet(cfg)
		assert.Error(t, err)
	})
}
