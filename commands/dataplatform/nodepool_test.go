package dataplatform

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCreateNodePoolRequest = resources.CreateNodePoolRequest{
		CreateNodePoolRequest: ionoscloud.CreateNodePoolRequest{
			Properties: &ionoscloud.CreateNodePoolProperties{
				Name:             &testNodePoolVar,
				NodeCount:        &testNodePoolIntVar,
				CpuFamily:        &testNodePoolVar,
				AvailabilityZone: &testAvailabilityZone,
				RamSize:          &testNodePoolIntVar,
				StorageSize:      &testNodePoolIntVar,
				StorageType:      &testStorageType,
				CoresCount:       &testNodePoolIntVar,
				Annotations:      &testNodePoolKVMap,
				Labels:           &testNodePoolKVMap,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
		},
	}
	testPatchNodePoolRequest = resources.PatchNodePoolRequest{
		PatchNodePoolRequest: ionoscloud.PatchNodePoolRequest{
			Properties: &ionoscloud.PatchNodePoolProperties{
				NodeCount: &testNodePoolIntNewVar,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolNewVar,
					Time:         &testNodePoolNewVar,
				},
				Annotations: &testNodePoolKVNewMap,
				Labels:      &testNodePoolKVNewMap,
			},
		},
	}
	testPatchOldNodePoolRequest = resources.PatchNodePoolRequest{
		PatchNodePoolRequest: ionoscloud.PatchNodePoolRequest{
			Properties: &ionoscloud.PatchNodePoolProperties{
				NodeCount: &testNodePoolIntVar,
			},
		},
	}
	testNodePoolGetNew = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				DataPlatformVersion: &testNodePoolNewVar,
				NodeCount:           &testNodePoolIntNewVar,
				Annotations:         &testNodePoolKVNewMap,
				Labels:              &testNodePoolKVNewMap,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolNewVar,
					Time:         &testNodePoolNewVar,
				}},
		},
	}
	testNodePoolGetOld = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DataPlatformVersion: &testNodePoolVar,
			},
		},
	}
	testNodePoolGet = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DatacenterId:        &testNodePoolVar,
				CpuFamily:           &testNodePoolVar,
				AvailabilityZone:    &testAvailabilityZone,
				RamSize:             &testNodePoolIntVar,
				StorageSize:         &testNodePoolIntVar,
				StorageType:         &testStorageType,
				DataPlatformVersion: &testNodePoolVar,
				CoresCount:          &testNodePoolIntVar,
				Annotations:         &testNodePoolKVMap,
				Labels:              &testNodePoolKVMap,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
			Metadata: &ionoscloud.Metadata{
				State: &testClusterStateVar,
			},
		},
	}
	nodePoolTestGet = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DatacenterId:        &testNodePoolVar,
				CpuFamily:           &testNodePoolVar,
				AvailabilityZone:    &testAvailabilityZone,
				RamSize:             &testNodePoolIntVar,
				StorageSize:         &testNodePoolIntVar,
				StorageType:         &testStorageType,
				DataPlatformVersion: &testNodePoolVar,
				CoresCount:          &testNodePoolIntVar,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
		},
	}
	nodePoolTestGetNew = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DatacenterId:        &testNodePoolVar,
				CpuFamily:           &testNodePoolVar,
				AvailabilityZone:    &testAvailabilityZone,
				RamSize:             &testNodePoolIntVar,
				StorageSize:         &testNodePoolIntVar,
				StorageType:         &testStorageType,
				DataPlatformVersion: &testNodePoolVar,
				CoresCount:          &testNodePoolIntVar,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolVar,
					Time:         &testNodePoolVar,
				},
			},
			Metadata: &ionoscloud.Metadata{
				State: &testNodePoolStateVar,
			},
		},
	}
	nodePoolTestGetFailed = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DatacenterId:        &testNodePoolVar,
				CpuFamily:           &testNodePoolVar,
				AvailabilityZone:    &testAvailabilityZone,
				RamSize:             &testNodePoolIntVar,
				StorageSize:         &testNodePoolIntVar,
				StorageType:         &testStorageType,
				DataPlatformVersion: &testNodePoolVar,
				CoresCount:          &testNodePoolIntVar,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolVar,
					Time:         &testNodePoolVar,
				},
			},
			Metadata: &ionoscloud.Metadata{
				State: &testNodePoolStateFailedVar,
			},
		},
	}

	nodePoolsList = resources.NodePoolListResponseData{
		NodePoolListResponseData: ionoscloud.NodePoolListResponseData{
			Id: &testNodePoolVar,
			Items: &[]ionoscloud.NodePoolResponseData{
				nodePoolTestId.NodePoolResponseData,
				nodePoolTestId.NodePoolResponseData,
			},
		},
	}
	nodePoolTestId = resources.NodePoolResponseData{
		NodePoolResponseData: ionoscloud.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &ionoscloud.NodePool{
				Name:                &testNodePoolVar,
				NodeCount:           &testNodePoolIntVar,
				DatacenterId:        &testNodePoolVar,
				CpuFamily:           &testNodePoolVar,
				AvailabilityZone:    &testAvailabilityZone,
				RamSize:             &testNodePoolIntVar,
				StorageSize:         &testNodePoolIntVar,
				StorageType:         &testStorageType,
				DataPlatformVersion: &testNodePoolVar,
				CoresCount:          &testNodePoolIntVar,
				MaintenanceWindow: &ionoscloud.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
			Metadata: &ionoscloud.Metadata{
				State: &testClusterStateVar,
			},
		},
	}

	nodePools = resources.NodePoolListResponseData{
		NodePoolListResponseData: ionoscloud.NodePoolListResponseData{
			Id:    &testNodePoolVar,
			Items: &[]ionoscloud.NodePoolResponseData{testNodePoolGet.NodePoolResponseData},
		},
	}

	testNodePoolKVMap          = map[string]interface{}{testNodePoolVar: testNodePoolVar}
	testNodePoolKVNewMap       = map[string]interface{}{testNodePoolNewVar: testNodePoolNewVar}
	testNodePoolIntVar         = int32(1)
	testNodePoolIntNewVar      = int32(1)
	testNodePoolVar            = "test-nodepool"
	testNodePoolNewVar         = "test-new-nodepool"
	testNodePoolStateFailedVar = "FAILED"
	testNodePoolStateVar       = "AVAILABLE"
	testNodePoolErr            = errors.New("nodepool test error")
	testAvailabilityZone       = ionoscloud.AvailabilityZone(testNodePoolVar)
	testStorageType            = ionoscloud.StorageType(testNodePoolVar)
)

func TestPreRunClusterNodePoolIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		err := PreRunClusterNodePoolIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterNodePoolIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunClusterNodePoolIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunNodePoolsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		err := PreRunNodePoolsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNodePoolsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunNodePoolsList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunNodePoolCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolVar)
		err := PreRunClusterNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNodePoolCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunClusterNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		err := PreRunClusterNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunClusterNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunNodePoolDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		err := PreRunClusterNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNodePoolDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		err := PreRunClusterNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(nodePools, nil, nil)
		err := RunNodePoolList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(nodePools, nil, testNodePoolErr)
		err := RunNodePoolList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		err := RunNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		err := RunNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGetFailed, nil, testNodePoolErr)
		err := RunNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, testNodePoolErr)
		err := RunNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCpuFamily), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAvailabilityZone), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgRam), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageType), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageSize), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCores), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(testNodePoolGet, nil, nil)
		err := RunNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolCreateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCpuFamily), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAvailabilityZone), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgRam), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageType), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageSize), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCores), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(nodePoolTestId, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestId, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestId, nil, nil)
		err := RunNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCpuFamily), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAvailabilityZone), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgRam), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageType), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageSize), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCores), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(nodePoolTestGetFailed, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGetFailed, nil, nil)
		err := RunNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCpuFamily), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAvailabilityZone), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgRam), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageType), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgStorageSize), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgCores), testNodePoolIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(testNodePoolGet, nil, testNodePoolErr)
		err := RunNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGetNew, nil, nil)
		err := RunNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolUpdateWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGetNew, nil, nil)
		err := RunNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGetNew, nil, testNodePoolErr)
		err := RunNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchOldNodePoolRequest).Return(testNodePoolGetOld, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		err := RunNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, testNodePoolErr)
		err := RunNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testNodePoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgAnnotations), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgLabels), testNodePoolKVNewMap)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(nodePoolTestGet, nil, testNodePoolErr)
		err := RunNodePoolUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(nodePoolsList, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(nodePoolsList, nil, testNodePoolErr)
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(resources.NodePoolListResponseData{}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(
			resources.NodePoolListResponseData{NodePoolListResponseData: ionoscloud.NodePoolListResponseData{Items: &[]ionoscloud.NodePoolResponseData{}}}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(nodePoolsList, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, testNodePoolErr)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, testNodePoolErr)
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, nil)
		err := RunNodePoolDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodePoolId), testNodePoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		cfg.Stdin = os.Stdin
		err := RunNodePoolDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNodePoolCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Name"})
	getNodePoolCols(core.GetGlobalFlagName("nodepool", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNodePoolColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nodepool", config.ArgCols), []string{"Unknown"})
	getNodePoolCols(core.GetGlobalFlagName("nodePool", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
