package pg

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-pg"
	"github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testBackup = resources.ClusterBackup{
		ClusterBackup: sdkgo.ClusterBackup{
			Id:          &testBackupVar,
			ClusterId:   &testBackupVar,
			DisplayName: &testBackupVar,
			Type:        &testBackupVar,
			Metadata: &sdkgo.Metadata{
				CreatedDate:      &testIonosTime,
				LastModifiedDate: &testIonosTime,
			},
		},
	}
	testBackups = resources.ClusterBackupList{
		ClusterBackupList: sdkgo.ClusterBackupList{
			Page: &testBackupPage,
			Data: &[]sdkgo.ClusterBackup{testBackup.ClusterBackup},
		},
	}
	testBackupPage = int32(2)
	testBackupVar  = "test-backup"
	testBackupErr  = errors.New("test backup error")
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunBackupId(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultBackupCols)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCols), defaultBackupCols)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgClusterId), testBackupVar)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgClusterId), testBackupVar)
		rm.CloudApiDbaasPgsqlMocks.Backup.EXPECT().ListBackups(testBackupVar).Return(testBackups, nil, testBackupErr)
		err := RunClusterBackupList(cfg)
		assert.Error(t, err)
	})
}

func TestGetBackupsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("backup", config.ArgCols), []string{"DisplayName"})
	getBackupCols(core.GetGlobalFlagName("backup", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetBackupsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("backup", config.ArgCols), []string{"Unknown"})
	getBackupCols(core.GetGlobalFlagName("backup", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
