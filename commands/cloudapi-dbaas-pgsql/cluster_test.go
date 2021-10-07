package cloudapi_dbaas_pgsql

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	v6resources "github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCreateClusterRequest = resources.CreateClusterRequest{
		CreateClusterRequest: sdkgo.CreateClusterRequest{
			DisplayName:     &testClusterVar,
			PostgresVersion: &testClusterVar,
			Location:        &testClusterVar,
			Replicas:        &testClusterFloatVar,
			RamSize:         &testClusterVar,
			CpuCoreCount:    &testClusterFloatVar,
			StorageSize:     &testClusterVar,
			StorageType:     &testClusterStorageTypeVar,
			VdcConnections: &[]sdkgo.VDCConnection{
				{
					VdcId:     &testClusterVar,
					LanId:     &testClusterVar,
					IpAddress: &testClusterVar,
				},
			},
			MaintenanceWindow: &sdkgo.MaintenanceWindow{
				Time:    &testClusterVar,
				Weekday: &testClusterVar,
			},
			Credentials: &sdkgo.DBUser{
				Username: &testClusterVar,
				Password: &testClusterVar,
			},
		},
	}
	testPatchClusterRequest = resources.PatchClusterRequest{
		PatchClusterRequest: sdkgo.PatchClusterRequest{
			CpuCoreCount: &testClusterFloatVar,
			RamSize:      &testClusterNewVar,
			StorageSize:  &testClusterNewVar,
			DisplayName:  &testClusterNewVar,
			MaintenanceWindow: &sdkgo.MaintenanceWindow{
				Time:    &testClusterNewVar,
				Weekday: &testClusterNewVar,
			},
			PostgresVersion: &testClusterNewVar,
		},
	}
	testCreateRestoreRequest = resources.CreateRestoreRequest{
		CreateRestoreRequest: sdkgo.CreateRestoreRequest{
			BackupId:           &testClusterVar,
			RecoveryTargetTime: &testClusterVar,
		},
	}
	testClusterGetNew = resources.Cluster{
		Cluster: sdkgo.Cluster{
			Id:              &testClusterVar,
			LifecycleStatus: &testClusterStateVar,
			DisplayName:     &testClusterNewVar,
			PostgresVersion: &testClusterNewVar,
			Location:        &testClusterVar,
			Replicas:        &testClusterFloatVar,
			RamSize:         &testClusterNewVar,
			CpuCoreCount:    &testClusterFloatVar,
			StorageSize:     &testClusterNewVar,
			BackupEnabled:   &testClusterBoolVar,
			StorageType:     &testClusterStorageTypeVar,
			VdcConnections: &[]sdkgo.VDCConnection{
				{
					VdcId:     &testClusterVar,
					LanId:     &testClusterVar,
					IpAddress: &testClusterVar,
				},
			},
			MaintenanceWindow: &sdkgo.MaintenanceWindow{
				Time:    &testClusterNewVar,
				Weekday: &testClusterNewVar,
			},
		},
	}
	testClusterGet = resources.Cluster{
		Cluster: sdkgo.Cluster{
			Id:              &testClusterVar,
			LifecycleStatus: &testClusterStateVar,
			DisplayName:     &testClusterVar,
			PostgresVersion: &testClusterVar,
			Location:        &testClusterVar,
			Replicas:        &testClusterFloatVar,
			RamSize:         &testClusterVar,
			CpuCoreCount:    &testClusterFloatVar,
			StorageSize:     &testClusterVar,
			BackupEnabled:   &testClusterBoolVar,
			StorageType:     &testClusterStorageTypeVar,
			VdcConnections: &[]sdkgo.VDCConnection{
				{
					VdcId:     &testClusterVar,
					LanId:     &testClusterVar,
					IpAddress: &testClusterVar,
				},
			},
			MaintenanceWindow: &sdkgo.MaintenanceWindow{
				Time:    &testClusterVar,
				Weekday: &testClusterVar,
			},
		},
	}
	testClusterGetFailed = resources.Cluster{
		Cluster: sdkgo.Cluster{
			Id:              &testClusterVar,
			LifecycleStatus: &testClusterStateFailedVar,
			DisplayName:     &testClusterVar,
			PostgresVersion: &testClusterVar,
			Location:        &testClusterVar,
			Replicas:        &testClusterFloatVar,
			RamSize:         &testClusterVar,
			CpuCoreCount:    &testClusterFloatVar,
			StorageSize:     &testClusterVar,
			StorageType:     &testClusterStorageTypeVar,
			VdcConnections: &[]sdkgo.VDCConnection{
				{
					VdcId:     &testClusterVar,
					LanId:     &testClusterVar,
					IpAddress: &testClusterVar,
				},
			},
			MaintenanceWindow: &sdkgo.MaintenanceWindow{
				Time:    &testClusterVar,
				Weekday: &testClusterVar,
			},
		},
	}
	testClusters = resources.ClusterList{
		ClusterList: sdkgo.ClusterList{
			Page: &testClusterPage,
			Data: &[]sdkgo.Cluster{testClusterGet.Cluster},
		},
	}
	testVdcGet = v6resources.Datacenter{
		Datacenter: ionoscloud.Datacenter{
			Id: &testClusterVar,
			Properties: &ionoscloud.DatacenterProperties{
				Location: &testClusterVar,
			},
		},
	}
	testClusterBoolVar        = true
	testClusterStateFailedVar = "FAILED"
	testClusterStateVar       = "AVAILABLE"
	testClusterFloatVar       = float32(1)
	testClusterStorageTypeVar = sdkgo.HDD
	testClusterPage           = int32(2)
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

func TestPreClusterId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		err := PreRunClusterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreClusterIdErr(t *testing.T) {
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

func TestPreRunClusterBackupIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		err := PreRunClusterBackupIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterBackupIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunClusterBackupIds(cfg)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		err := PreRunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterCreateRecoveryTime(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		err := PreRunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterCreateRecoveryTimeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		err := PreRunClusterCreate(cfg)
		assert.Error(t, err)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultClusterCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
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
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List("").Return(testClusters, nil, testClusterErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, testClusterErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGetFailed, nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest, testClusterVar, testClusterVar).Return(&testClusterGet, nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest, testClusterVar, testClusterVar).Return(&testClusterGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateLocation(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testClusterVar).Return(&testVdcGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest, testClusterVar, testClusterVar).Return(&testClusterGet, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateLocationErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testClusterVar).Return(&testVdcGet, nil, testClusterErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest, testClusterVar, testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil).Times(2)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgVdcId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgIpAddress), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgReplicas), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageType), string(sdkgo.HDD))
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest, testClusterVar, testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGetFailed, nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterNewVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(&testClusterGetNew, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterRestore(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Restore.EXPECT().Restore(testClusterVar, testCreateRestoreRequest).Return(nil, nil)
		err := RunClusterRestore(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterRestoreErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Restore.EXPECT().Restore(testClusterVar, testCreateRestoreRequest).Return(nil, testClusterErr)
		err := RunClusterRestore(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterRestoreAskConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgTime), testClusterVar)
		cfg.Stdin = os.Stdin
		err := RunClusterRestore(cfg)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, nil)
		err := RunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, testClusterErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, testClusterErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		cfg.Stdin = os.Stdin
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgPostgresVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgRamSize), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgCpuCoreCount), testClusterFloatVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStorageSize), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceDay), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgName), testClusterNewVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(&testClusterGetNew, nil, testClusterErr)
		err := RunClusterUpdate(cfg)
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
