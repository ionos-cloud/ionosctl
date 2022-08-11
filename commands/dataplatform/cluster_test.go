package dataplatform

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCreateClusterRequest = resources.CreateClusterRequest{
		CreateClusterRequest: sdkgo.CreateClusterRequest{
			Properties: &sdkgo.CreateClusterProperties{
				DatacenterId:        &testClusterVar,
				Name:                &testClusterVar,
				DataPlatformVersion: &testClusterVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testClusterVar,
					DayOfTheWeek: &testClusterVar,
				},
			},
		},
	}
	testPatchClusterRequest = resources.PatchClusterRequest{
		PatchClusterRequest: sdkgo.PatchClusterRequest{
			Properties: &sdkgo.PatchClusterProperties{
				Name:                &testClusterNewVar,
				DataPlatformVersion: &testClusterNewVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testClusterNewVar,
					DayOfTheWeek: &testClusterNewVar,
				},
			},
		},
	}
	testPatchClusterRequestOld = resources.PatchClusterRequest{
		PatchClusterRequest: sdkgo.PatchClusterRequest{
			Properties: &sdkgo.PatchClusterProperties{
				Name: &testClusterVar,
			},
		},
	}
	testClusterGetNew = resources.ClusterResponseData{
		ClusterResponseData: sdkgo.ClusterResponseData{
			Id: &testClusterVar,
			Properties: &sdkgo.Cluster{
				DatacenterId:        &testClusterNewVar,
				Name:                &testClusterNewVar,
				DataPlatformVersion: &testClusterNewVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testClusterNewVar,
					DayOfTheWeek: &testClusterNewVar,
				},
			},
			Metadata: &sdkgo.Metadata{
				State: &testClusterStateVar,
			},
		},
	}
	testClusterGet = resources.ClusterResponseData{
		ClusterResponseData: sdkgo.ClusterResponseData{
			Id: &testClusterVar,
			Properties: &sdkgo.Cluster{
				DatacenterId:        &testClusterVar,
				Name:                &testClusterVar,
				DataPlatformVersion: &testClusterVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testClusterVar,
					DayOfTheWeek: &testClusterVar,
				},
			},
			Metadata: &sdkgo.Metadata{
				State: &testClusterStateVar,
			},
		},
	}
	testClusterGetFailed = resources.ClusterResponseData{
		ClusterResponseData: sdkgo.ClusterResponseData{
			Id: &testClusterVar,
			Properties: &sdkgo.Cluster{
				DatacenterId:        &testClusterVar,
				Name:                &testClusterVar,
				DataPlatformVersion: &testClusterVar,
				MaintenanceWindow: &sdkgo.MaintenanceWindow{
					Time:         &testClusterVar,
					DayOfTheWeek: &testClusterVar,
				},
			},
			Metadata: &sdkgo.Metadata{
				State: &testClusterStateFailedVar,
			},
		},
	}
	testClusters = resources.ClusterListResponseData{
		ClusterListResponseData: sdkgo.ClusterListResponseData{
			Items: &[]sdkgo.ClusterResponseData{testClusterGet.ClusterResponseData},
		},
	}

	testClusterStateFailedVar = "FAILED"
	testClusterStateVar       = "AVAILABLE"
	testClusterVar            = "test-cluster"
	testClusterNewVar         = "test-cluster-new"
	testClusterErr            = errors.New("test cluster error")
)

func TestClusterCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ClusterCmd())
	if ok := ClusterCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunClusterId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		err := PreRunClusterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunClusterId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		err := PreRunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		err := PreRunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterDeleteNameErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		err := PreRunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultClusterCols)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		err := RunClusterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.DataPlatformMocks.Cluster.EXPECT().List("").Return(testClusters, nil, testClusterErr)
		err := RunClusterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, nil)
		err := RunClusterGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, testClusterErr)
		err := RunClusterGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterGetWaitForState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, nil)
		err := RunClusterGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterGetWaitForStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetFailed, nil, nil)
		err := RunClusterGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(testClusterGet, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(testClusterGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(testClusterGet, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, nil).Times(2)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateWaitForStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(testClusterGet, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetFailed, nil, nil)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForStateIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(resources.ClusterResponseData{}, nil, nil)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForStateNewClusterErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(testClusterGet, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterNewVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(testClusterGetNew, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetNew, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterNewVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(testClusterGetNew, nil, testClusterErr)
		err := RunClusterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequestOld).Return(testClusterGetNew, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetNew, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateWaitForRequest(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterNewVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetNew, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(testClusterGetNew, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetNew, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateWaitForRequestErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgMaintenanceDay), testClusterNewVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Get(testClusterVar).Return(testClusterGetFailed, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(testClusterGetFailed, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Delete(testClusterVar).Return(resources.ClusterResponseData{}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Delete(testClusterVar).Return(resources.ClusterResponseData{}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, testClusterErr)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(resources.ClusterListResponseData{}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(
			resources.ClusterListResponseData{ClusterListResponseData: sdkgo.ClusterListResponseData{Items: &[]sdkgo.ClusterResponseData{}}}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgName), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		rm.DataPlatformMocks.Cluster.EXPECT().Delete(testClusterVar).Return(resources.ClusterResponseData{}, nil, testClusterErr)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		rm.DataPlatformMocks.Cluster.EXPECT().Delete(testClusterVar).Return(resources.ClusterResponseData{}, nil, testClusterErr)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAskConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.DataPlatformMocks.Cluster.EXPECT().Delete(testClusterVar).Return(resources.ClusterResponseData{}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterDeleteAskConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dp.ArgClusterId), testClusterVar)
		cfg.Stdin = os.Stdin
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetClustersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("cluster", config.ArgCols), []string{"DisplayName"})
	getClusterCols(core.GetGlobalFlagName("cluster", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetClustersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("cluster", config.ArgCols), []string{"Unknown"})
	getClusterCols(core.GetGlobalFlagName("cluster", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
