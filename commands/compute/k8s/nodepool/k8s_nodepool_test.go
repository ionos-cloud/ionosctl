package nodepool

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	k8scluster "github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testK8sNodePoolLanBoolVar = false
	nodepoolTestPost          = resources.K8sNodePoolForPost{
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
				State: &testutil.TestStateVar,
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
				State: &testutil.TestStateVar,
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
				State: &testutil.TestStateVar,
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
	testClusterVar        = "test-cluster"
	clusterTestId         = resources.K8sCluster{
		KubernetesCluster: ionoscloud.KubernetesCluster{
			Id: &testClusterVar,
			Properties: &ionoscloud.KubernetesClusterProperties{
				Name:       &testClusterVar,
				K8sVersion: &testClusterVar,
				S3Buckets: &[]ionoscloud.S3Bucket{
					{
						Name: &testClusterVar,
					},
				},
				ApiSubnetAllowList: &[]string{testClusterVar},
			},
		},
	}
	clusterTest = resources.K8sCluster{
		KubernetesCluster: ionoscloud.KubernetesCluster{
			Properties: &ionoscloud.KubernetesClusterProperties{
				Name:       &testClusterVar,
				K8sVersion: &testClusterVar,
				S3Buckets: &[]ionoscloud.S3Bucket{
					{
						Name: &testClusterVar,
					},
				},
				ApiSubnetAllowList: &[]string{testClusterVar},
			},
		},
	}
	clustersList = resources.K8sClusters{
		KubernetesClusters: ionoscloud.KubernetesClusters{
			Id: &testClusterVar,
			Items: &[]ionoscloud.KubernetesCluster{
				clusterTestId.KubernetesCluster,
				clusterTestId.KubernetesCluster,
			},
		},
	}
)

func TestPreRunK8sClusterNodePoolIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testutil.TestQueryParamVar))
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListClusters().Return(clustersList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testClusterVar).Return(nodepoolsList, &testutil.TestResponse, nil).Times(len(k8scluster.GetK8sClusters(clustersList)))
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, &testutil.TestResponse, nil)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		cfg.SetFlag(constants.FlagOrderBy, testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(
			resources.K8sNodePools{
				KubernetesNodePools: ionoscloud.KubernetesNodePools{
					Items: &[]ionoscloud.KubernetesNodePool{}},
			},
			&testutil.TestResponse, nil)
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
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepools, nil, testNodepoolErr)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgWait, true)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgName, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCpuFamily, testNodepoolVar)
		cfg.SetFlag(constants.FlagAvailabilityZone, testNodepoolVar)
		cfg.SetFlag(constants.FlagRam, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagStorageType, testNodepoolVar)
		cfg.SetFlag(constants.FlagStorageSize, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCores, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgK8sVersion, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgLanIds, []int{int(testNodepoolIntVar)})
		cfg.SetFlag(cloudapiv6.ArgDhcp, testK8sNodePoolLanBoolVar)
		cfg.SetFlag(constants.FlagLabels, testNodepoolKVMap)
		cfg.SetFlag(constants.FlagAnnotations, testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgName, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCpuFamily, testNodepoolVar)
		cfg.SetFlag(constants.FlagAvailabilityZone, testNodepoolVar)
		cfg.SetFlag(constants.FlagRam, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagStorageType, testNodepoolVar)
		cfg.SetFlag(constants.FlagStorageSize, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCores, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetCluster(testNodepoolVar).Return(&clusterTest, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, gomock.Any()).Return(&nodepoolTestPrivate, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgName, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCpuFamily, testNodepoolVar)
		cfg.SetFlag(constants.FlagAvailabilityZone, testNodepoolVar)
		cfg.SetFlag(constants.FlagRam, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagStorageType, testNodepoolVar)
		cfg.SetFlag(constants.FlagStorageSize, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCores, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgLanIds, []int{int(testNodepoolIntVar)})
		cfg.SetFlag(cloudapiv6.ArgDhcp, testK8sNodePoolLanBoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetCluster(testNodepoolVar).Return(nil, nil, testNodepoolErr)
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
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgName, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCpuFamily, testNodepoolVar)
		cfg.SetFlag(constants.FlagAvailabilityZone, testNodepoolVar)
		cfg.SetFlag(constants.FlagRam, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagStorageType, testNodepoolVar)
		cfg.SetFlag(constants.FlagStorageSize, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagCores, testNodepoolIntVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgK8sVersion, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgLanIds, []int{int(testNodepoolIntVar)})
		cfg.SetFlag(cloudapiv6.ArgDhcp, testK8sNodePoolLanBoolVar)
		cfg.SetFlag(constants.FlagLabels, testNodepoolKVMap)
		cfg.SetFlag(constants.FlagAnnotations, testNodepoolKVMap)
		rm.CloudApiV6Mocks.K8s.EXPECT().CreateNodePool(testNodepoolVar, nodepoolTestPost).Return(&nodepoolTest, nil, testNodepoolErr)
		err := RunK8sNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodePoolUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgK8sVersion, testNodepoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgK8sMaintenanceDay, testNodepoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgK8sMaintenanceTime, testNodepoolNewVar)
		cfg.SetFlag(constants.FlagAnnotations, testNodepoolKVNewMap)
		cfg.SetFlag(constants.FlagLabels, testNodepoolKVNewMap)
		cfg.SetFlag(cloudapiv6.ArgK8sMinNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(cloudapiv6.ArgK8sMaxNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(cloudapiv6.ArgLanIds, []int{int(testNodepoolIntNewVar)})
		cfg.SetFlag(cloudapiv6.ArgDhcp, testK8sNodePoolLanBoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgPublicIps, []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().UpdateNodePool(testNodepoolVar, testNodepoolVar, nodepoolTestUpdateNew).Return(&nodepoolTestNew, nil, testNodepoolErr)
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
		cfg.SetFlag(cloudapiv6.ArgK8sVersion, testNodepoolNewVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgK8sMaintenanceDay, testNodepoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgK8sMaintenanceTime, testNodepoolNewVar)
		cfg.SetFlag(constants.FlagAnnotations, testNodepoolKVNewMap)
		cfg.SetFlag(constants.FlagLabels, testNodepoolKVNewMap)
		cfg.SetFlag(cloudapiv6.ArgK8sMinNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(cloudapiv6.ArgK8sMaxNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(constants.FlagNodeCount, testNodepoolIntNewVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testNodepoolIntNewVar)
		cfg.SetFlag(cloudapiv6.ArgDhcp, testK8sNodePoolLanBoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(cloudapiv6.ArgPublicIps, []string{testNodepoolNewVar, testNodepoolNewVar})
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNodePool(testNodepoolVar, testNodepoolVar).Return(&nodepoolTestGet, nil, testNodepoolErr)
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
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepoolsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepoolsList, nil, testNodepoolErr)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(resources.K8sNodePools{}, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(
			resources.K8sNodePools{KubernetesNodePools: ionoscloud.KubernetesNodePools{Items: &[]ionoscloud.KubernetesNodePool{}}}, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodePools(testNodepoolVar).Return(nodepoolsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testutil.TestResponse, testNodepoolErr)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, testNodepoolErr)
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
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNodePool(testNodepoolVar, testNodepoolVar).Return(nil, nil)
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
		cfg.SetFlag(constants.FlagNodepoolId, testNodepoolVar)
		cfg.SetFlag(constants.FlagClusterId, testNodepoolVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunK8sNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}
