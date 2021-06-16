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
	nodepoolTestPost = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:             &testNodepoolVar,
				NodeCount:        &testNodepoolIntVar,
				DatacenterId:     &testNodepoolVar,
				CpuFamily:        &testNodepoolVar,
				AvailabilityZone: &testNodepoolVar,
				RamSize:          &testNodepoolIntVar,
				StorageSize:      &testNodepoolIntVar,
				StorageType:      &testNodepoolVar,
				K8sVersion:       &testNodepoolVar,
				CoresCount:       &testNodepoolIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTest = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:             &testNodepoolVar,
				NodeCount:        &testNodepoolIntVar,
				DatacenterId:     &testNodepoolVar,
				CpuFamily:        &testNodepoolVar,
				AvailabilityZone: &testNodepoolVar,
				RamSize:          &testNodepoolIntVar,
				StorageSize:      &testNodepoolIntVar,
				StorageType:      &testNodepoolVar,
				K8sVersion:       &testNodepoolVar,
				CoresCount:       &testNodepoolIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTestId = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:             &testNodepoolVar,
				NodeCount:        &testNodepoolIntVar,
				DatacenterId:     &testNodepoolVar,
				CpuFamily:        &testNodepoolVar,
				AvailabilityZone: &testNodepoolVar,
				RamSize:          &testNodepoolIntVar,
				StorageSize:      &testNodepoolIntVar,
				StorageType:      &testNodepoolVar,
				K8sVersion:       &testNodepoolVar,
				CoresCount:       &testNodepoolIntVar,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepoolTestGet = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:                     &testNodepoolVar,
				NodeCount:                &testNodepoolIntVar,
				DatacenterId:             &testNodepoolVar,
				CpuFamily:                &testNodepoolVar,
				AvailabilityZone:         &testNodepoolVar,
				RamSize:                  &testNodepoolIntVar,
				StorageSize:              &testNodepoolIntVar,
				StorageType:              &testNodepoolVar,
				K8sVersion:               &testNodepoolVar,
				CoresCount:               &testNodepoolIntVar,
				PublicIps:                &testNodepoolSliceVar,
				AvailableUpgradeVersions: &testNodepoolSliceVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolVar,
					Time:         &testNodepoolVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntVar,
					MaxNodeCount: &testNodepoolIntVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepoolTestGetNew = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:                     &testNodepoolVar,
				NodeCount:                &testNodepoolIntVar,
				DatacenterId:             &testNodepoolVar,
				CpuFamily:                &testNodepoolVar,
				AvailabilityZone:         &testNodepoolVar,
				RamSize:                  &testNodepoolIntVar,
				StorageSize:              &testNodepoolIntVar,
				StorageType:              &testNodepoolVar,
				K8sVersion:               &testNodepoolVar,
				CoresCount:               &testNodepoolIntVar,
				PublicIps:                &testNodepoolSliceVar,
				AvailableUpgradeVersions: &testNodepoolSliceVar,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolVar,
					Time:         &testNodepoolVar,
				},
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntVar,
					MaxNodeCount: &testNodepoolIntVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testNodepoolIntNewVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepools = resources.K8sNodePools{
		KubernetesNodePools: ionoscloud.KubernetesNodePools{
			Id:    &testNodepoolVar,
			Items: &[]ionoscloud.KubernetesNodePool{nodepoolTest.KubernetesNodePool},
		},
	}
	nodepoolTestNew = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				Name:       &testNodepoolVar,
				K8sVersion: &testNodepoolNewVar,
				NodeCount:  &testNodepoolIntNewVar,
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				Annotations: &nodepoolTestMap,
				Labels:      &nodepoolTestMap,
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testNodepoolIntNewVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
				PublicIps: &[]string{testNodepoolNewVar, testNodepoolNewVar},
			},
		},
	}
	nodepoolTestMap = map[string]string{
		testNodepoolNewVar: testNodepoolNewVar,
	}
	nodepoolTestUpdateNew = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				K8sVersion: &testNodepoolNewVar,
				NodeCount:  &testNodepoolIntNewVar,
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Annotations: &nodepoolTestMap,
				Labels:      &nodepoolTestMap,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
					{
						Id:   &testNodepoolIntNewVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
				PublicIps: &[]string{testNodepoolNewVar, testNodepoolNewVar},
			},
		},
	}
	nodepoolTestOld = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				Name:       &testNodepoolVar,
				NodeCount:  &testNodepoolIntVar,
				K8sVersion: &testNodepoolVar,
			},
		},
	}
	nodepoolTestUpdateOld = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				NodeCount:  &testNodepoolIntVar,
				K8sVersion: &testNodepoolVar,
			},
		},
	}
	testNodepoolIntVar    = int32(1)
	testNodepoolIntNewVar = int32(1)
	testNodepoolVar       = "test-nodepool"
	testNodepoolSliceVar  = []string{"test-nodepool"}
	testNodepoolNewVar    = "test-new-nodepool"
	testNodepoolErr       = errors.New("nodepool test error")
)

func TestPreRunK8sClusterNodePoolIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodePoolIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sClusterDcIdsNodePoolName(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		err := PreRunK8sClusterDcIdsNodePoolName(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterDcIdsNodePoolNameErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterDcIdsNodePoolName(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, nil, nil)
		err := RunK8sNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, nil, testNodepoolErr)
		err := RunK8sNodePoolList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		err := RunK8sNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		err := RunK8sNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, nil, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolCreateGetK8sVersionErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().GetVersion().Return(testNodepoolVar, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateWaitStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGetNew, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGetNew, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGetNew, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateOld).Return(&nodepoolTestOld, nil, nil)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, testNodepoolErr)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgK8sClusterId), testNodepoolVar)
		cfg.Stdin = os.Stdin
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetK8sNodePoolCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Name"})
	getK8sNodePoolCols(core.GetGlobalFlagName("nodepool", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodePoolColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Unknown"})
	getK8sNodePoolCols(core.GetGlobalFlagName("nodepool", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetK8sNodePoolsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getK8sNodePoolsIds(w, testNodepoolVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
