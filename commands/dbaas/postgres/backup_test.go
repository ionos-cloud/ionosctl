package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testBackup = resources.BackupResponse{
		BackupResponse: sdkgo.BackupResponse{
			Id: &testBackupVar,
			Properties: &sdkgo.ClusterBackup{
				Id:                         &testBackupVar,
				ClusterId:                  &testBackupVar,
				EarliestRecoveryTargetTime: &testIonosTime,
				Version:                    &testBackupVar,
				IsActive:                   &testBackupBoolVar,
			},
			Metadata: &sdkgo.BackupMetadata{
				State:       &testStateVar,
				CreatedDate: &testIonosTime,
			},
		},
	}
	testBackups = resources.ClusterBackupList{
		ClusterBackupList: sdkgo.ClusterBackupList{
			Id:    &testBackupVar,
			Items: &[]sdkgo.BackupResponse{testBackup.BackupResponse},
		},
	}
	testStateVar      = sdkgo.State("AVAILABLE")
	testBackupVar     = "test-backup"
	testBackupBoolVar = true
	testBackupErr     = errors.New("test backup error")
)

func TestBackupCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(BackupCmd())
	if ok := BackupCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestClusterBackupCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ClusterBackupCmd())
	if ok := ClusterBackupCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreBackupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testBackupVar)
		err := PreRunBackupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreBackupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunBackupId(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultBackupCols)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().List().Return(testBackups, nil, nil)
		err := RunBackupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().List().Return(testBackups, nil, testBackupErr)
		err := RunBackupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testBackupVar)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().Get(testBackupVar).Return(&testBackup, nil, nil)
		err := RunBackupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgBackupId), testBackupVar)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().Get(testBackupVar).Return(&testBackup, nil, testBackupErr)
		err := RunBackupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterBackupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultBackupCols)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testBackupVar)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().ListBackups(testBackupVar).Return(testBackups, nil, nil)
		err := RunClusterBackupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterBackupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testBackupVar)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().ListBackups(testBackupVar).Return(testBackups, nil, testBackupErr)
		err := RunClusterBackupList(cfg)
		assert.Error(t, err)
	})
}
