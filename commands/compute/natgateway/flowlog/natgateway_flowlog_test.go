package flowlog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testNatGatewayFlowLogErr = errors.New("natgateway-rule test error")

func TestPreRunNatGatewayFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		err := PreRunNATGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayFlowLogListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testutil.TestQueryParamVar))
		err := PreRunNATGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		err := PreRunNATGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogVar)
		err := PreRunNatGatewayFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		err := PreRunDcNatGatewayFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcNatGatewayFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(testutil.TestFlowLogs, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagOrderBy), testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(resources.FlowLogs{}, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(testutil.TestFlowLogs, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestFlowLog, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestFlowLog, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestInputFlowLog).Return(&testutil.TestFlowLog, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestInputFlowLog).Return(&testutil.TestFlowLog, &testutil.TestResponse, testutil.TestFlowLogErr)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestInputFlowLog).Return(&testutil.TestFlowLog, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestInputFlowLog).Return(&testutil.TestFlowLog, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar, &testutil.TestFlowLogProperties).Return(&testutil.TestFlowLogUpdated, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar, &testutil.TestFlowLogProperties).Return(&testutil.TestFlowLogUpdated, nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testutil.TestFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testutil.TestFlowLogNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar, &testutil.TestFlowLogProperties).Return(&testutil.TestFlowLogUpdated, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(testutil.TestFlowLogsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(testutil.TestFlowLogsList, nil, testutil.TestFlowLogErr)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(resources.FlowLogs{}, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(
			resources.FlowLogs{FlowLogs: ionoscloud.FlowLogs{Items: &[]ionoscloud.FlowLog{}}}, &testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListFlowLogs(testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(testutil.TestFlowLogsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, testutil.TestFlowLogErr)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(nil, testNatGatewayFlowLogErr)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteFlowLog(testutil.TestFlowLogVar, testutil.TestFlowLogVar, testutil.TestFlowLogVar).Return(nil, nil)
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testutil.TestFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testutil.TestFlowLogVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNatGatewayFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}
