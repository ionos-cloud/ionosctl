package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testFlowLog = v6.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Id: &testFlowLogVar,
			Properties: &ionoscloud.FlowLogProperties{
				Name:      &testFlowLogVar,
				Action:    &testFlowLogUpperVar,
				Direction: &testFlowLogUpperVar,
				Bucket:    &testFlowLogVar,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testFlowLogState,
			},
		},
	}
	testInputFlowLog = v6.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: testFlowLog.FlowLog.Properties,
		},
	}
	testFlowLogUpdated = v6.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: &testFlowLogProperties.FlowLogProperties,
		},
	}
	testFlowLogProperties = v6.FlowLogProperties{
		FlowLogProperties: ionoscloud.FlowLogProperties{
			Name:      &testFlowLogNewVar,
			Action:    &testFlowLogNewUpperVar,
			Direction: &testFlowLogNewUpperVar,
			Bucket:    &testFlowLogNewVar,
		},
	}
	testFlowLogs = v6.FlowLogs{
		FlowLogs: ionoscloud.FlowLogs{
			Id:    &testFlowLogVar,
			Items: &[]ionoscloud.FlowLog{testFlowLog.FlowLog},
		},
	}
	testFlowLogState       = "AVAILABLE"
	testFlowLogVar         = "test-flowlog"
	testFlowLogUpperVar    = strings.ToUpper(testFlowLogVar)
	testFlowLogNewVar      = "test-new-flowlog"
	testFlowLogNewUpperVar = strings.ToUpper(testFlowLogNewVar)
	testFlowLogErr         = errors.New("flowlog test error")
)

func TestPreRunFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Bucket), testFlowLogVar)
		err := PreRunFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFlowLogId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		err := PreRunDcServerNicFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFlowLogIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerNicFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		rm.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, nil)
		err := RunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		rm.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, testFlowLogErr)
		err := RunFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.FlowLog.EXPECT().Get(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, testFlowLogErr)
		err := RunFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.FlowLog.EXPECT().Get(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, nil)
		err := RunFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Bucket), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FlowLog.EXPECT().Create(testFlowLogVar, testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testInputFlowLog, nil, nil)
		err := RunFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Bucket), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FlowLog.EXPECT().Create(testFlowLogVar, testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testInputFlowLog, &testResponse, nil)
		err := RunFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, testFlowLogErr)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetFlowLogsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Name"})
	getFlowLogsCols(core.GetGlobalFlagName("firewallrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetFlowLogsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Unknown"})
	getFlowLogsCols(core.GetGlobalFlagName("firewallrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetFlowLogsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getFlowLogsIds(w, testFlowLogVar, testFlowLogVar, testFlowLogVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
