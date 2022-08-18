package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
				Annotations:      &testNodepoolKVMap,
				Labels:           &testNodepoolKVMap,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTestPrivatePost = resources.K8sNodePoolForPost{
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
				CoresCount:       &testNodepoolIntVar,
				K8sVersion:       &testNodepoolVar,
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
				Annotations:      &testNodepoolKVMap,
				Labels:           &testNodepoolKVMap,
				Lans: &[]ionoscloud.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTestPrivate = resources.K8sNodePool{
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
	nodepoolsList = resources.K8sNodePools{
		KubernetesNodePools: ionoscloud.KubernetesNodePools{
			Id: &testNodepoolVar,
			Items: &[]ionoscloud.KubernetesNodePool{
				nodepoolTestId.KubernetesNodePool,
				nodepoolTestId.KubernetesNodePool,
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
				Annotations: &testNodepoolKVNewMap,
				Labels:      &testNodepoolKVNewMap,
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
				Annotations: &testNodepoolKVNewMap,
				Labels:      &testNodepoolKVNewMap,
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
	testNodepoolKVMap     = map[string]string{testNodepoolVar: testNodepoolVar}
	testNodepoolKVNewMap  = map[string]string{testNodepoolNewVar: testNodepoolNewVar}
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
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

func TestPreRunK8sNodePoolsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		err := PreRunK8sClusterNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		err := PreRunK8sClusterNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sNodePoolsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunK8sNodePoolsList(cfg)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
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

func TestRunK8sNodePoolListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListClusters(gomock.AssignableToTypeOf(testListQueryParam)).Return(clustersList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testClusterVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepoolsList, &testResponse, nil).Times(len(getK8sClusters(clustersList)))
		err := RunK8sNodePoolListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepools, &testResponse, nil)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.K8sNodePools{}, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepools, nil, testNodepoolErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTest, &testResponse, nil)
		err := RunK8sNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolPrivateCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetVersion().Return(testNodepoolVar, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPrivatePost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestPrivate, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetVersion().Return(testNodepoolVar, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestId, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTest, nil, testNodepoolErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestNew, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGetNew, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestNew, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGetNew, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestNew, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGetNew, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateOld, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestOld, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestNew, nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublicIps), []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepoolsList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepoolsList, nil, testNodepoolErr)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.K8sNodePools{}, &testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.K8sNodePools{KubernetesNodePools: ionoscloud.KubernetesNodePools{Items: &[]ionoscloud.KubernetesNodePool{}}}, &testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodepoolsList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testNodepoolErr)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testNodepoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodepoolVar)
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
	getK8sNodePoolCols(core.GetGlobalFlagName("nodepool", config.ArgCols), core.GetFlagName("nodepool", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodePoolColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Unknown"})
	getK8sNodePoolCols(core.GetGlobalFlagName("nodepool", config.ArgCols), core.GetFlagName("nodepool", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
