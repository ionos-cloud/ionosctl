package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
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
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testNodepoolVar,
			},
		},
	}
	nodepools = resources.K8sNodePools{
		KubernetesNodePools: ionoscloud.KubernetesNodePools{
			Id:    &testNodepoolVar,
			Items: &[]ionoscloud.KubernetesNodePool{nodepoolTest.KubernetesNodePool},
		},
	}
	nodepoolTestNew = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				K8sVersion: &testNodepoolNewVar,
				NodeCount:  &testNodepoolIntNewVar,
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				Annotations: &ionoscloud.KubernetesNodePoolAnnotation{
					Key:   &testNodepoolNewVar,
					Value: &testNodepoolNewVar,
				},
				Labels: &ionoscloud.KubernetesNodePoolLabel{
					Key:   &testNodepoolNewVar,
					Value: &testNodepoolNewVar,
				},
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id: &testNodepoolIntNewVar,
					},
				},
			},
		},
	}
	nodepoolTestUpdatedNew = resources.K8sNodePoolUpdated{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
				Name:             &testNodepoolVar,
				DatacenterId:     &testNodepoolVar,
				CpuFamily:        &testNodepoolVar,
				AvailabilityZone: &testNodepoolVar,
				StorageType:      &testNodepoolVar,
				K8sVersion:       &testNodepoolNewVar,
				NodeCount:        &testNodepoolIntNewVar,
				RamSize:          &testNodepoolIntVar,
				StorageSize:      &testNodepoolIntVar,
				CoresCount:       &testNodepoolIntVar,
				AutoScaling: &ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				Annotations: &ionoscloud.KubernetesNodePoolAnnotation{
					Key:   &testNodepoolNewVar,
					Value: &testNodepoolNewVar,
				},
				Labels: &ionoscloud.KubernetesNodePoolLabel{
					Key:   &testNodepoolNewVar,
					Value: &testNodepoolNewVar,
				},
				MaintenanceWindow: &ionoscloud.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id: &testNodepoolIntNewVar,
					},
				},
			},
		},
	}
	nodepoolTestOld = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				NodeCount:  &testNodepoolIntVar,
				K8sVersion: &testNodepoolVar,
			},
		},
	}
	nodepoolTestUpdatedOld = resources.K8sNodePoolUpdated{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
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
			},
		},
	}
	testNodepoolIntVar    = int32(1)
	testNodepoolIntNewVar = int32(2)
	testNodepoolVar       = "test-nodepool"
	testNodepoolSliceVar  = []string{"test-nodepool"}
	testNodepoolNewVar    = "test-new-nodepool"
	testNodepoolErr       = errors.New("nodepool test error")
)

func TestPreRunK8sClusterNodePoolIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodePoolIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), "")
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sClusterDcIdsNodePoolName(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolName), testNodepoolVar)
		err := PreRunK8sClusterDcIdsNodePoolName(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterDcIdsNodePoolNameErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), "")
		err := PreRunK8sClusterDcIdsNodePoolName(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, nil, nil)
		err := RunK8sNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, nil, testNodepoolErr)
		err := RunK8sNodePoolList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		err := RunK8sNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolName), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeZone), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRamSize), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgStorageType), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgCoresCount), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolVersion), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testNodepoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTest).Return(&nodepoolTest, nil, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolName), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgCpuFamily), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeZone), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRamSize), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgStorageType), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgStorageSize), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgCoresCount), testNodepoolIntVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolVersion), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testNodepoolVar)
		rm.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTest).Return(&nodepoolTest, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolVersion), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestNew).Return(&nodepoolTestUpdatedNew, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestOld).Return(&nodepoolTestUpdatedOld, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolVersion), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestNew).Return(&nodepoolTestUpdatedNew, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolVersion), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testNodepoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testNodepoolIntNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		rm.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, testNodepoolErr)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodepoolVar)
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
	viper.Set(builder.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Name"})
	getK8sNodePoolCols(builder.GetGlobalFlagName("nodepool", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodePoolColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Unknown"})
	getK8sNodePoolCols(builder.GetGlobalFlagName("nodepool", config.ArgCols), w)
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
