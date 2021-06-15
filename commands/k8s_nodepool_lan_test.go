package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	k8sNodepoolLanTest = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				NodeCount: &testK8sNodePoolLanIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testK8sNodePoolLanIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	inputK8sNodepoolLanTest = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				NodeCount: &testK8sNodePoolLanIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testK8sNodePoolLanIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testK8sNodePoolLanNewIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	k8sNodepoolLanTestUpdated = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testK8sNodePoolLanVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				NodeCount: &testK8sNodePoolLanIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testK8sNodePoolLanIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testK8sNodePoolLanNewIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	testK8sNodePoolLanIntVar    = int32(1)
	testK8sNodePoolLanNewIntVar = int32(2)
	testK8sNodePoolLanBoolVar   = true
	testK8sNodePoolLanVar       = "test-nodepool-lan"
	testK8sNodePoolLanErr       = errors.New("nodepool-lan test error")
)

func TestPreRunK8sClusterNodePoolLanIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanVar)
		err := PreRunK8sClusterNodePoolLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodePoolLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterNodePoolLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgCols), defaultK8sNodePoolLanCols)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		err := RunK8sNodePoolLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanNewIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar, inputK8sNodepoolLanTest).Return(&k8sNodepoolLanTestUpdated, nil, nil)
		err := RunK8sNodePoolLanAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolLanAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanNewIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar, inputK8sNodepoolLanTest).Return(&k8sNodepoolLanTestUpdated, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanNewIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestGetK8sNodePoolLanCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"LanId"})
	getK8sNodePoolLanCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodePoolLanColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"Unknown"})
	getK8sNodePoolLanCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
