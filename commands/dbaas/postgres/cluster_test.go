package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCreateClusterRequest = resources.CreateClusterRequest{
		CreateClusterRequest: psql.CreateClusterRequest{
			Properties: &psql.CreateClusterProperties{
				DisplayName:         testClusterVar,
				PostgresVersion:     testClusterVar,
				Location:            testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				Instances:           testClusterIntVar,
				Ram:                 testClusterIntVar,
				Cores:               testClusterIntVar,
				StorageSize:         testClusterIntVar,
				SynchronizationMode: testSyncModeVar,
				StorageType:         psql.StorageType(testClusterStorageTypeVar),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
				Credentials: psql.DBUser{
					Username: testClusterVar,
					Password: testClusterVar,
				},
				FromBackup: &psql.CreateRestoreRequest{
					BackupId:           testClusterVar,
					RecoveryTargetTime: &testIonosTime,
				},
			},
		},
	}
	testCreateClusterRequestSSDPremium = resources.CreateClusterRequest{
		CreateClusterRequest: psql.CreateClusterRequest{
			Properties: &psql.CreateClusterProperties{
				DisplayName:         testClusterVar,
				PostgresVersion:     testClusterVar,
				Location:            testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				Instances:           testClusterIntVar,
				Ram:                 testClusterIntVar,
				Cores:               testClusterIntVar,
				StorageSize:         testClusterIntVar,
				SynchronizationMode: testSyncModeVar,
				StorageType:         psql.StorageType(testClusterStorageTypeSSDPremiumVar),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
				Credentials: psql.DBUser{
					Username: testClusterVar,
					Password: testClusterVar,
				},
				FromBackup: &psql.CreateRestoreRequest{
					BackupId:           testClusterVar,
					RecoveryTargetTime: &testIonosTime,
				},
			},
		},
	}
	testCreateClusterRequestSSDStandard = resources.CreateClusterRequest{
		CreateClusterRequest: psql.CreateClusterRequest{
			Properties: &psql.CreateClusterProperties{
				DisplayName:         testClusterVar,
				PostgresVersion:     testClusterVar,
				Location:            testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				Instances:           testClusterIntVar,
				Ram:                 testClusterIntVar,
				Cores:               testClusterIntVar,
				StorageSize:         testClusterIntVar,
				SynchronizationMode: testSyncModeVar,
				StorageType:         psql.StorageType(testClusterStorageTypeSSDStandardVar),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
				Credentials: psql.DBUser{
					Username: testClusterVar,
					Password: testClusterVar,
				},
				FromBackup: &psql.CreateRestoreRequest{
					BackupId:           testClusterVar,
					RecoveryTargetTime: &testIonosTime,
				},
			},
		},
	}
	testPatchClusterRequest = resources.PatchClusterRequest{
		PatchClusterRequest: psql.PatchClusterRequest{
			Properties: &psql.PatchClusterProperties{
				Cores:       &testClusterIntNewVar,
				Ram:         &testClusterIntNewVar,
				StorageSize: &testClusterIntNewVar,
				DisplayName: &testClusterNewVar,
				Instances:   &testClusterIntNewVar,
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterNewVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterNewVar),
				},
				PostgresVersion: &testClusterNewVar,
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterNewVar,
						LanId:        testClusterNewVar,
						Cidr:         testClusterNewVar,
					},
				},
			},
		},
	}
	testCreateRestoreRequest = resources.CreateRestoreRequest{
		CreateRestoreRequest: psql.CreateRestoreRequest{
			BackupId:           testClusterVar,
			RecoveryTargetTime: &testIonosTime,
		},
	}
	testClusterGetNew = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterNewVar,
				PostgresVersion:     &testClusterNewVar,
				SynchronizationMode: &testSyncModeVar,
				Location:            &testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntNewVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntNewVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeVar)),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterNewVar,
						LanId:        testClusterNewVar,
						Cidr:         testClusterNewVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterNewVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterNewVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateVar),
			},
		},
	}
	testClusterGet = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterVar,
				PostgresVersion:     &testClusterVar,
				Location:            &testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				SynchronizationMode: &testSyncModeVar,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeVar)),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateVar),
			},
		},
	}
	testClusterGetSSDPremium = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterVar,
				PostgresVersion:     &testClusterVar,
				Location:            &testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				SynchronizationMode: &testSyncModeVar,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeSSDPremiumVar)),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateVar),
			},
		},
	}
	testClusterGetSSDStandard = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterVar,
				PostgresVersion:     &testClusterVar,
				Location:            &testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				SynchronizationMode: &testSyncModeVar,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeSSDStandardVar)),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateVar),
			},
		},
	}
	testClusterGetNoConnection = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterVar,
				PostgresVersion:     &testClusterVar,
				Location:            &testClusterVar,
				SynchronizationMode: &testSyncModeVar,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeVar)),
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateVar),
			},
		},
	}
	testClusterGetFailed = resources.ClusterResponse{
		ClusterResponse: psql.ClusterResponse{
			Id: &testClusterVar,
			Properties: &psql.ClusterProperties{
				DisplayName:         &testClusterVar,
				SynchronizationMode: &testSyncModeVar,
				PostgresVersion:     &testClusterVar,
				Location:            &testClusterVar,
				BackupLocation:      &testClusterBackupLocation,
				Instances:           &testClusterIntVar,
				Ram:                 &testClusterIntVar,
				Cores:               &testClusterIntVar,
				StorageSize:         &testClusterIntVar,
				StorageType:         pointer.From(psql.StorageType(testClusterStorageTypeVar)),
				Connections: []psql.Connection{
					{
						DatacenterId: testClusterVar,
						LanId:        testClusterVar,
						Cidr:         testClusterVar,
					},
				},
				MaintenanceWindow: &psql.MaintenanceWindow{
					Time:         testClusterVar,
					DayOfTheWeek: psql.DayOfTheWeek(testClusterVar),
				},
			},
			Metadata: &psql.ClusterMetadata{
				State: (*psql.State)(&testClusterStateFailedVar),
			},
		},
	}
	testClusters = resources.ClusterList{
		ClusterList: psql.ClusterList{
			Items: []psql.ClusterResponse{testClusterGet.ClusterResponse},
		},
	}
	testVdcGet = cloudapiv6resources.Datacenter{
		Datacenter: compute.Datacenter{
			Id: &testClusterVar,
			Properties: &compute.DatacenterProperties{
				Location: &testClusterVar,
			},
		},
	}
	testIonosTime = psql.IonosTime{
		Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	testSyncModeVar                      = psql.SynchronizationMode(strings.ToUpper(testClusterVar))
	testTimeArgVar                       = "2021-01-01T00:00:00Z"
	testClusterBoolVar                   = true
	testClusterStateFailedVar            = "FAILED"
	testClusterStateVar                  = "AVAILABLE"
	testClusterMbVar                     = "1MB"
	testClusterBackupLocation            = "de"
	testClusterIntVar                    = int32(1)
	testClusterIntNewVar                 = int32(2)
	testClusterStorageTypeVar            = "HDD"
	testClusterStorageTypeSSDPremiumVar  = "SSD Premium"
	testClusterStorageTypeSSDStandardVar = "SSD Standard"
	testClusterVar                       = "test-cluster"
	testClusterNewVar                    = "test-cluster-new"
	testClusterErr                       = errors.New("test cluster error")
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		err := PreRunClusterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreClusterIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunClusterId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterBackupIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		err := PreRunClusterBackupIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterBackupIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunClusterBackupIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		err := PreRunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterCreateInstancesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), 10)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		err := PreRunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterCreateCoresErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), 0)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		err := PreRunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterCreateRecoveryTime(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		err := PreRunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		err := PreRunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterDeleteNameErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		err := PreRunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultClusterCols)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateSSDPremium(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), "SSD_Premium")
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequestSSDPremium).Return(&testClusterGetSSDPremium, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateSSDStandard(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), "SSD_standard")
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequestSSDStandard).Return(&testClusterGetSSDStandard, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateConvertErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), "test")
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateLocation(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testClusterVar, gomock.AssignableToTypeOf(cloudapiv6resources.QueryParams{})).Return(&testVdcGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, nil)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateLocationErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), strings.ToUpper(testClusterVar))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testClusterVar, gomock.AssignableToTypeOf(cloudapiv6resources.QueryParams{})).Return(&testVdcGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil).Times(2)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateWaitForStateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, &resources.Response{}, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil).Times(2)
		err := RunClusterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterCreateWaitForStateIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&resources.ClusterResponse{}, nil, nil)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGetFailed, nil, nil)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterCreateWaitForStateNewClusterErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), testClusterBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLocation), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSyncMode), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterMbVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageType), string("HDD"))
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbUsername), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDbPassword), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupLocation), testClusterBackupLocation)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Create(testCreateClusterRequest).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, testClusterErr)
		err := RunClusterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterNewVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(&testClusterGetNew, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateRemoveConnection(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRemoveConnection), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, resources.PatchClusterRequest{PatchClusterRequest: psql.PatchClusterRequest{
			Properties: &psql.PatchClusterProperties{
				Connections: []psql.Connection{},
			},
		}}).Return(&testClusterGetNoConnection, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateRemoveConnectionErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRemoveConnection), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, testClusterErr)
		err := RunClusterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdateRemoveConnectionAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRemoveConnection), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, resources.PatchClusterRequest{PatchClusterRequest: psql.PatchClusterRequest{
			Properties: &psql.PatchClusterProperties{
				Connections: []psql.Connection{},
			},
		}}).Return(&testClusterGetNoConnection, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterUpdateRemoveConnectionAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRemoveConnection), true)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, resources.PatchClusterRequest{PatchClusterRequest: psql.PatchClusterRequest{
			Properties: &psql.PatchClusterProperties{
				Connections: []psql.Connection{},
			},
		}}).Return(&testClusterGetNoConnection, nil, nil)
		err := RunClusterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterRestore(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRecoveryTime), testTimeArgVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunClusterRestore(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Delete(testClusterVar).Return(nil, nil)
		err := RunClusterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(testClusters, nil, testClusterErr)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(resources.ClusterList{}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().List(testClusterVar).Return(
			resources.ClusterList{ClusterList: psql.ClusterList{Items: []psql.ClusterResponse{}}}, nil, nil)
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunClusterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgVersion), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgRam), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgInstances), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCores), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStorageSize), testClusterIntNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceDay), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgMaintenanceTime), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgName), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDatacenterId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLanId), testClusterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgCidr), testClusterNewVar)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Get(testClusterVar).Return(&testClusterGet, nil, nil)
		rm.CloudApiDbaasPgsqlMocks.Cluster.EXPECT().Update(testClusterVar, testPatchClusterRequest).Return(&testClusterGetNew, nil, testClusterErr)
		err := RunClusterUpdate(cfg)
		assert.Error(t, err)
	})
}
