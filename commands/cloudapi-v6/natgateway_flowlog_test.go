package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testNatGatewayFlowLogErr = errors.New("natgateway-rule test error")

func TestPreRunNatGatewayFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogVar)
		err := PreRunNatGatewayFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		err := PreRunDcNatGatewayFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNatGatewayFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, nil)
		err := RunNatGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, nil)
		err := RunNatGatewayFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, nil)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, &testResponse, nil)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, nil)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, nil)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, nil)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgFlowLogId), testFlowLogVar)
		cfg.Stdin = os.Stdin
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}
