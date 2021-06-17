package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
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
				NodeCount:  &testK8sNodePoolLanIntVar,
				K8sVersion: &testK8sNodePoolLanVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testK8sNodePoolLanVar,
					Time:         &testK8sNodePoolLanVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testK8sNodePoolLanIntVar,
					MaxNodeCount: &testK8sNodePoolLanIntVar,
				},
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
				NodeCount:  &testK8sNodePoolLanIntVar,
				K8sVersion: &testK8sNodePoolLanVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testK8sNodePoolLanVar,
					Time:         &testK8sNodePoolLanVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testK8sNodePoolLanIntVar,
					MaxNodeCount: &testK8sNodePoolLanIntVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testK8sNodePoolLanIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testK8sNodePoolLanNewIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
						Routes: &[]ionoscloud.KubernetesNodePoolLanRoutes{
							{
								Network:   &testK8sNodePoolLanVar,
								GatewayIp: &testK8sNodePoolLanVar,
							},
						},
					},
				},
			},
		},
	}
	inputK8sNodepoolLanTestRemove = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				NodeCount:  &testK8sNodePoolLanIntVar,
				K8sVersion: &testK8sNodePoolLanVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testK8sNodePoolLanVar,
					Time:         &testK8sNodePoolLanVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testK8sNodePoolLanIntVar,
					MaxNodeCount: &testK8sNodePoolLanIntVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{},
			},
		},
	}
	k8sNodepoolLanTestUpdatedRemove = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testK8sNodePoolLanVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				NodeCount:  &testK8sNodePoolLanIntVar,
				K8sVersion: &testK8sNodePoolLanVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testK8sNodePoolLanVar,
					Time:         &testK8sNodePoolLanVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testK8sNodePoolLanIntVar,
					MaxNodeCount: &testK8sNodePoolLanIntVar,
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	k8sNodepoolLanTestUpdated = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testK8sNodePoolLanVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				NodeCount:  &testK8sNodePoolLanIntVar,
				K8sVersion: &testK8sNodePoolLanVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testK8sNodePoolLanVar,
					Time:         &testK8sNodePoolLanVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testK8sNodePoolLanIntVar,
					MaxNodeCount: &testK8sNodePoolLanIntVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testK8sNodePoolLanIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testK8sNodePoolLanNewIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
						Routes: &[]ionoscloud.KubernetesNodePoolLanRoutes{
							{
								Network:   &testK8sNodePoolLanVar,
								GatewayIp: &testK8sNodePoolLanVar,
							},
						},
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

func TestRunK8sNodePoolLanListLansErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&resources.K8sNodePool{
			KubernetesNodePool: ionoscloud.KubernetesNodePool{
				Id:         &testK8sNodePoolLanVar,
				Properties: &ionoscloud.KubernetesNodePoolProperties{},
			},
		}, nil, nil)
		err := RunK8sNodePoolLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanListPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&resources.K8sNodePool{
			KubernetesNodePool: ionoscloud.KubernetesNodePool{
				Id: &testK8sNodePoolLanVar,
			},
		}, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetwork), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGatewayIp), testK8sNodePoolLanVar)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetwork), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGatewayIp), testK8sNodePoolLanVar)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetwork), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGatewayIp), testK8sNodePoolLanVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar, inputK8sNodepoolLanTestRemove).Return(&k8sNodepoolLanTestUpdatedRemove, nil, nil)
		err := RunK8sNodePoolLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolLanRemoveAsk(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar, inputK8sNodepoolLanTestRemove).Return(&k8sNodepoolLanTestUpdatedRemove, nil, nil)
		err := RunK8sNodePoolLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolLanRemoveAskErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanIntVar)
		err := RunK8sNodePoolLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar, inputK8sNodepoolLanTestRemove).Return(&k8sNodepoolLanTestUpdatedRemove, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolLanRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testK8sNodePoolLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testK8sNodePoolLanIntVar)
		rm.K8s.EXPECT().GetNodePool(testK8sNodePoolLanVar, testK8sNodePoolLanVar).Return(&k8sNodepoolLanTest, nil, testK8sNodePoolLanErr)
		err := RunK8sNodePoolLanRemove(cfg)
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
