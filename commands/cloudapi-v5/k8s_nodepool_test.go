package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	nodepoolTestPost = resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: ionoscloud.KubernetesNodePoolForPost{
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPost{
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
						Id: &testNodepoolIntNewVar,
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
	nodepoolsList = resources.K8sNodePools{
		KubernetesNodePools: ionoscloud.KubernetesNodePools{
			Id: &testNodepoolVar,
			Items: &[]ionoscloud.KubernetesNodePool{
				nodepoolTestId.KubernetesNodePool,
				nodepoolTestId.KubernetesNodePool,
			},
		},
	}
	nodepoolTestNew = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolProperties{
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
						Id: &testNodepoolIntNewVar,
					},
				},
				PublicIps: &[]string{testNodepoolNewVar, testNodepoolNewVar},
			},
		},
	}
	nodepoolTestMap = map[string]string{
		testNodepoolNewVar: testNodepoolNewVar,
	}
	nodepoolTestUpdateNew = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
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
						Id: &testNodepoolIntNewVar,
					},
				},
				PublicIps: &[]string{testNodepoolNewVar, testNodepoolNewVar},
			},
		},
	}
	nodepoolTestOld = resources.K8sNodePool{
		KubernetesNodePool: ionoscloud.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &ionoscloud.KubernetesNodePoolProperties{
				Name:       &testNodepoolVar,
				NodeCount:  &testNodepoolIntVar,
				K8sVersion: &testNodepoolVar,
			},
		},
	}
	nodepoolTestUpdateOld = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &ionoscloud.KubernetesNodePoolPropertiesForPut{
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

func TestK8sNodePoolCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(K8sNodePoolCmd())
	if ok := K8sNodePoolCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunK8sNodePoolsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolsListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunK8sNodePoolsList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sClusterNodePoolIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
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

func TestPreRunK8sClusterDcIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		err := PreRunK8sClusterDcIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterDcIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterDcIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCols), allK8sNodePoolCols)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, resources.ListQueryParams{}).Return(nodepools, &testResponse, nil)
		err := RunK8sNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, testListQueryParam).Return(resources.K8sNodePools{}, &testResponse, nil)
		err := RunK8sNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, resources.ListQueryParams{}).Return(nodepools, nil, testNodepoolErr)
		err := RunK8sNodePoolList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, &testResponse, nil)
		err := RunK8sNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		err := RunK8sNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, &testResponse, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolCreateGetK8sVersionErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetVersion().Return(testNodepoolVar, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateWaitStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestId, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, &testResponse, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateOld).Return(&nodepoolTestOld, nil, nil)
		err := RunK8sNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sAnnotationValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV5Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
		err := RunK8sNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, resources.ListQueryParams{}).Return(nodepoolsList, &testResponse, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV5Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, testNodepoolErr)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV5Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgK8sClusterId), testNodepoolVar)
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
