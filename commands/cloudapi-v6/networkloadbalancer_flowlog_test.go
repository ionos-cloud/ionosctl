package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testNetworkLoadBalancerFlowLogErr = errors.New("networkloadbalancer-rule test error")

func TestPreRunNetworkLoadBalancerFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		err := PreRunNetworkLoadBalacerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerFlowLogListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunNetworkLoadBalacerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunNetworkLoadBalacerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		err := PreRunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		err := PreRunDcNetworkLoadBalancerFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagOrderBy), testQueryParamVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(resources.FlowLogs{}, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogs, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testFlowLog, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, &testResponse, testFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateFlowLog(testFlowLogVar, testFlowLogVar, testInputFlowLog).Return(&testFlowLog, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).
			Return(&testFlowLogUpdated, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties).Return(&testFlowLogUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerFlowLogUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogsList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogsList, nil, testFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(resources.FlowLogs{}, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(
			resources.FlowLogs{FlowLogs: ionoscloud.FlowLogs{Items: &[]ionoscloud.FlowLog{}}}, &testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListFlowLogs(testFlowLogVar, testFlowLogVar).Return(testFlowLogsList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, testFlowLogErr)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, testNetworkLoadBalancerFlowLogErr)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteFlowLog(testFlowLogVar, testFlowLogVar, testFlowLogVar).Return(nil, nil)
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNetworkLoadBalancerFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}
