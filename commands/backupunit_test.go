package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	backupUnitTest = v5.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Email:    &testBackupUnitVar,
				Name:     &testBackupUnitVar,
				Password: &testBackupUnitVar,
			},
		},
	}
	backupUnitTestGet = v5.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Id:         &testBackupUnitVar,
			Properties: backupUnitTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	backupUnitTestGetSSO = v5.BackupUnitSSO{
		BackupUnitSSO: ionoscloud.BackupUnitSSO{
			SsoUrl: &testBackupUnitVar,
		},
	}
	backupUnits = v5.BackupUnits{
		BackupUnits: ionoscloud.BackupUnits{
			Id:    &testBackupUnitVar,
			Items: &[]ionoscloud.BackupUnit{backupUnitTest.BackupUnit},
		},
	}
	backupUnitProperties = v5.BackupUnitProperties{
		BackupUnitProperties: ionoscloud.BackupUnitProperties{
			Email:    &testBackupUnitNewVar,
			Password: &testBackupUnitNewVar,
		},
	}
	backupUnitNew = v5.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Name:     &testBackupUnitVar,
				Email:    &testBackupUnitNewVar,
				Password: &testBackupUnitNewVar,
			},
		},
	}
	testBackupUnitVar    = "test-backup-unit"
	testBackupUnitNewVar = "test-new-backup-unit"
	testBackupUnitErr    = errors.New("backup-unit test error")
)

func TestPreRunBackupUnitId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		err := PreRunBackupUnitId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunBackupUnitId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunBackupUnitNameEmailPwd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitVar)
		err := PreRunBackupUnitNameEmailPwd(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitNameEmailPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunBackupUnitNameEmailPwd(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.BackupUnit.EXPECT().List().Return(backupUnits, &testResponse, nil)
		err := RunBackupUnitList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.BackupUnit.EXPECT().List().Return(backupUnits, nil, testBackupUnitErr)
		err := RunBackupUnitList(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Get(testBackupUnitVar).Return(&backupUnitTestGet, &testResponse, nil)
		err := RunBackupUnitGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Get(testBackupUnitVar).Return(&backupUnitTestGet, nil, testBackupUnitErr)
		err := RunBackupUnitGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitGetSsoUrl(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().GetSsoUrl(testBackupUnitVar).Return(&backupUnitTestGetSSO, &testResponse, nil)
		err := RunBackupUnitGetSsoUrl(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitGetSsoUrlErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().GetSsoUrl(testBackupUnitVar).Return(&backupUnitTestGetSSO, nil, testBackupUnitErr)
		err := RunBackupUnitGetSsoUrl(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Create(backupUnitTest).Return(&backupUnitTest, &testResponse, nil)
		err := RunBackupUnitCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Create(backupUnitTest).Return(&backupUnitTest, &testResponseErr, nil)
		err := RunBackupUnitCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Create(backupUnitTest).Return(&backupUnitTest, nil, nil)
		err := RunBackupUnitCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Create(backupUnitTest).Return(&backupUnitTest, nil, testBackupUnitErr)
		err := RunBackupUnitCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitNewVar)
		rm.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties).Return(&backupUnitNew, &testResponse, nil)
		err := RunBackupUnitUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitNewVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties).Return(&backupUnitNew, nil, nil)
		err := RunBackupUnitUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testBackupUnitNewVar)
		rm.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties).Return(&backupUnitNew, nil, testBackupUnitErr)
		err := RunBackupUnitUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Delete(testBackupUnitVar).Return(&testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		rm.BackupUnit.EXPECT().Delete(testBackupUnitVar).Return(nil, testBackupUnitErr)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.BackupUnit.EXPECT().Delete(testBackupUnitVar).Return(nil, nil)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.BackupUnit.EXPECT().Delete(testBackupUnitVar).Return(nil, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBackupUnitId), testBackupUnitVar)
		cfg.Stdin = os.Stdin
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetBackupUnitCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("backupunit", config.ArgCols), []string{"Name"})
	getBackupUnitCols(core.GetGlobalFlagName("backupunit", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetBackupUnitColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("backupunit", config.ArgCols), []string{"Unknown"})
	getBackupUnitCols(core.GetGlobalFlagName("backupunit", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetBackupUnitIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getBackupUnitsIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
