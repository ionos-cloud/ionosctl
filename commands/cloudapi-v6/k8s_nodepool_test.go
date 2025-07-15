package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	nodepoolTestPost = resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: compute.KubernetesNodePoolForPost{
			Properties: &compute.KubernetesNodePoolPropertiesForPost{
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
				Lans: &[]compute.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTestPrivatePost = resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: compute.KubernetesNodePoolForPost{
			Properties: &compute.KubernetesNodePoolPropertiesForPost{
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
		KubernetesNodePool: compute.KubernetesNodePool{
			Properties: &compute.KubernetesNodePoolProperties{
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
				Lans: &[]compute.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
		},
	}
	nodepoolTestPrivate = resources.K8sNodePool{
		KubernetesNodePool: compute.KubernetesNodePool{
			Properties: &compute.KubernetesNodePoolProperties{
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
		KubernetesNodePools: compute.KubernetesNodePools{
			Id: &testNodepoolVar,
			Items: &[]compute.KubernetesNodePool{
				nodepoolTestId.KubernetesNodePool,
				nodepoolTestId.KubernetesNodePool,
			},
		},
	}
	nodepoolTestId = resources.K8sNodePool{
		KubernetesNodePool: compute.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &compute.KubernetesNodePoolProperties{
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
				Lans: &[]compute.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &compute.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepoolTestGet = resources.K8sNodePool{
		KubernetesNodePool: compute.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &compute.KubernetesNodePoolProperties{
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
				MaintenanceWindow: &compute.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolVar,
					Time:         &testNodepoolVar,
				},
				AutoScaling: &compute.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntVar,
					MaxNodeCount: &testNodepoolIntVar,
				},
				Lans: &[]compute.KubernetesNodePoolLan{
					{
						Id:   &testNodepoolIntVar,
						Dhcp: &testK8sNodePoolLanBoolVar,
					},
				},
			},
			Metadata: &compute.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepoolTestGetNew = resources.K8sNodePool{
		KubernetesNodePool: compute.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &compute.KubernetesNodePoolProperties{
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
				MaintenanceWindow: &compute.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolVar,
					Time:         &testNodepoolVar,
				},
				AutoScaling: &compute.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntVar,
					MaxNodeCount: &testNodepoolIntVar,
				},
				Lans: &[]compute.KubernetesNodePoolLan{
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
			Metadata: &compute.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	nodepools = resources.K8sNodePools{
		KubernetesNodePools: compute.KubernetesNodePools{
			Id:    &testNodepoolVar,
			Items: &[]compute.KubernetesNodePool{nodepoolTest.KubernetesNodePool},
		},
	}
	nodepoolTestNew = resources.K8sNodePool{
		KubernetesNodePool: compute.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &compute.KubernetesNodePoolProperties{
				Name:       &testNodepoolVar,
				K8sVersion: &testNodepoolNewVar,
				NodeCount:  &testNodepoolIntNewVar,
				AutoScaling: &compute.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				Annotations: &testNodepoolKVNewMap,
				Labels:      &testNodepoolKVNewMap,
				MaintenanceWindow: &compute.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Lans: &[]compute.KubernetesNodePoolLan{
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
		KubernetesNodePoolForPut: compute.KubernetesNodePoolForPut{
			Properties: &compute.KubernetesNodePoolPropertiesForPut{
				K8sVersion: &testNodepoolNewVar,
				NodeCount:  &testNodepoolIntNewVar,
				AutoScaling: &compute.KubernetesAutoScaling{
					MinNodeCount: &testNodepoolIntNewVar,
					MaxNodeCount: &testNodepoolIntNewVar,
				},
				MaintenanceWindow: &compute.KubernetesMaintenanceWindow{
					DayOfTheWeek: &testNodepoolNewVar,
					Time:         &testNodepoolNewVar,
				},
				Annotations: &testNodepoolKVNewMap,
				Labels:      &testNodepoolKVNewMap,
				Lans: &[]compute.KubernetesNodePoolLan{
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
		KubernetesNodePool: compute.KubernetesNodePool{
			Id: &testNodepoolVar,
			Properties: &compute.KubernetesNodePoolProperties{
				Name:       &testNodepoolVar,
				NodeCount:  &testNodepoolIntVar,
				K8sVersion: &testNodepoolVar,
			},
		},
	}
	nodepoolTestUpdateOld = resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: compute.KubernetesNodePoolForPut{
			Properties: &compute.KubernetesNodePoolPropertiesForPut{
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodePoolIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunK8sClusterNodePoolIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sNodePoolsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolsListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		err := PreRunK8sClusterNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodePoolDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunK8sClusterNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sNodePoolsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunK8sNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.K8sNodePools{
				KubernetesNodePools: compute.KubernetesNodePools{
					Items: &[]compute.KubernetesNodePool{}},
			},
			&testResponse, nil)
		err := RunK8sNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodePoolListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVMap)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetVersion().Return(testNodepoolVar, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolCreateWaitStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVMap)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageSize), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testNodepoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVMap)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanIds), []int{int(testNodepoolIntNewVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sVersion), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceDay), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaintenanceTime), testNodepoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAnnotations), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagLabels), testNodepoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMinNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sMaxNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodeCount), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNodepoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), testK8sNodePoolLanBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.K8sNodePools{KubernetesNodePools: compute.KubernetesNodePools{Items: &[]compute.KubernetesNodePool{}}}, &testResponse, nil)
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodepoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodepoolVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}
