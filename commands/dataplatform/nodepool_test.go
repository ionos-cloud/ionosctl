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
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCreateNodePoolRequest = resources.CreateNodePoolRequest{
		CreateNodePoolRequest: sdkgo.CreateNodePoolRequest{
			Properties: &sdkgo.CreateNodePoolProperties{
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
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
		},
	}
	testPatchNodePoolRequest = resources.PatchNodePoolRequest{
		PatchNodePoolRequest: sdkgo.PatchNodePoolRequest{
			Properties: &sdkgo.PatchNodePoolProperties{
				NodeCount: &testNodePoolIntNewVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolNewVar,
					Time:         &testNodePoolNewVar,
				},
				Annotations: &testNodePoolKVNewMap,
				Labels:      &testNodePoolKVNewMap,
			},
		},
	}
	testPatchOldNodePoolRequest = resources.PatchNodePoolRequest{
		PatchNodePoolRequest: sdkgo.PatchNodePoolRequest{
			Properties: &sdkgo.PatchNodePoolProperties{
				NodeCount: &testNodePoolIntVar,
			},
		},
	}
	testNodePoolGetNew = resources.NodePoolResponseData{
		NodePoolResponseData: sdkgo.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &sdkgo.NodePool{
				Name:                &testNodePoolVar,
				DataPlatformVersion: &testNodePoolNewVar,
				NodeCount:           &testNodePoolIntNewVar,
				Annotations:         &testNodePoolKVNewMap,
				Labels:              &testNodePoolKVNewMap,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolNewVar,
					Time:         &testNodePoolNewVar,
				}},
		},
	}
	testNodePoolGet = resources.NodePoolResponseData{
		NodePoolResponseData: sdkgo.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &sdkgo.NodePool{
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
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testNodePoolVar,
					DayOfTheWeek: &testNodePoolVar,
				},
			},
			Metadata: &sdkgo.Metadata{
				State: &testNodePoolStateVar,
			},
		},
	}
	testNodePoolGetFailed = resources.NodePoolResponseData{
		NodePoolResponseData: sdkgo.NodePoolResponseData{
			Id: &testNodePoolVar,
			Properties: &sdkgo.NodePool{
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
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					DayOfTheWeek: &testNodePoolVar,
					Time:         &testNodePoolVar,
				},
			},
			Metadata: &sdkgo.Metadata{
				State: &testNodePoolStateFailedVar,
			},
		},
	}

	testNodePools = resources.NodePoolListResponseData{
		NodePoolListResponseData: sdkgo.NodePoolListResponseData{
			Items: &[]sdkgo.NodePoolResponseData{testNodePoolGet.NodePoolResponseData},
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
	testAvailabilityZone       = sdkgo.AvailabilityZone(testNodePoolVar)
	testStorageType            = sdkgo.StorageType(testNodePoolVar)
)

func TestNodePoolCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(NodePoolCmd())
	if ok := NodePoolCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

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

func TestRunNodePoolList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testNodePoolVar)
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(testNodePools, nil, nil)
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
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(testNodePools, nil, testNodePoolErr)
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, nil)
		err := RunNodePoolGet(cfg)
		assert.NoError(t, err)
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, testNodePoolErr)
		err := RunNodePoolGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolGetWaitForState(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, &resources.Response{}, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, &resources.Response{}, nil)
		err := RunNodePoolGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolGetWaitForStateErr(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGetFailed, nil, testNodePoolErr)
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

func TestRunNodePoolCreateWaitForState(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(testNodePoolGet, &resources.Response{}, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, &resources.Response{}, nil).Times(2)
		err := RunNodePoolCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolCreateWaitForStateErr(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(testNodePoolGetFailed, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGetFailed, nil, nil)
		err := RunNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolCreateWaitForStateIdErr(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(resources.NodePoolResponseData{}, nil, nil)
		err := RunNodePoolCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNodePoolCreateWaitForStateNewNodePoolErr(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Create(testNodePoolVar, testCreateNodePoolRequest).Return(testNodePoolGet, &resources.Response{}, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, &resources.Response{}, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, &resources.Response{}, testNodePoolErr)
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
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, nil)
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
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, testNodePoolErr)
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
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgNodeCount), testNodePoolIntVar)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchOldNodePoolRequest).Return(testNodePoolGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, nil)
		err := RunNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolUpdateWaitForRequest(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetNew, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGet, nil, nil)
		err := RunNodePoolUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNodePoolUpdateWaitForRequestErr(t *testing.T) {
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
		rm.DataPlatformMocks.NodePool.EXPECT().Get(testNodePoolVar, testNodePoolVar).Return(testNodePoolGetFailed, nil, testNodePoolErr)
		rm.DataPlatformMocks.NodePool.EXPECT().Update(testNodePoolVar, testNodePoolVar, testPatchNodePoolRequest).Return(testNodePoolGetFailed, nil, nil)
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
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(testNodePools, nil, nil)
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
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(testNodePools, nil, testNodePoolErr)
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
			resources.NodePoolListResponseData{NodePoolListResponseData: sdkgo.NodePoolListResponseData{Items: &[]sdkgo.NodePoolResponseData{}}}, nil, nil)
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
		rm.DataPlatformMocks.NodePool.EXPECT().List(testNodePoolVar).Return(testNodePools, nil, nil)
		rm.DataPlatformMocks.NodePool.EXPECT().Delete(testNodePoolVar, testNodePoolVar).Return(resources.NodePoolResponseData{}, nil, testNodePoolErr)
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
