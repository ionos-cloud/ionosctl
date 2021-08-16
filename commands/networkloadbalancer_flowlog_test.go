package commands

import (
	"bufio"
	"bytes"
	"errors"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testNetworkLoadBalancerFlowLogErr = errors.New("networkloadbalancer-rule test error")

func TestPreRunNetworkLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogVar)
		err := PreRunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		err := PreRunDcNetworkLoadBalancerFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, nil)
		err := RunNetworkLoadBalancerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, nil)
		err := RunNetworkLoadBalancerFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, nil)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, nil)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogNewVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, nil)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogNewVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgBucketName), testFlowLogNewVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, nil)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		rm.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFlowLogId), testFlowLogVar)
		cfg.Stdin = os.Stdin
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNetworkLoadBalancerFlowLogsIds(t *testing.T) {
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
	getNetworkLoadBalancerFlowLogsIds(w, testFlowLogVar, testFlowLogVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
