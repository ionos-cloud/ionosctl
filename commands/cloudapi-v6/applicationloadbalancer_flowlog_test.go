package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testApplicationLoadBalancerFlowLogErr = errors.New("applicationloadbalancer-rule test error")

func TestApplicationLoadBalancerFlowLogCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ApplicationLoadBalancerFlowLogCmd())
	if ok := ApplicationLoadBalancerFlowLogCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunApplicationLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
			err := PreRunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunApplicationLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			err := PreRunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestPreRunApplicationLoadBalancerFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			err := PreRunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunApplicationLoadBalancerFlowLogDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			err := PreRunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunApplicationLoadBalancerFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			err := PreRunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestPreRunDcApplicationLoadBalancerFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			err := PreRunDcApplicationLoadBalancerFlowLogIds(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunDcApplicationLoadBalancerFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			err := PreRunDcApplicationLoadBalancerFlowLogIds(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, nil, nil)
			err := RunApplicationLoadBalancerFlowLogList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(
				core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters),
				[]string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)},
			)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, nil, nil)
			err := RunApplicationLoadBalancerFlowLogList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, &testResponse, nil)
			err := RunApplicationLoadBalancerFlowLogList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, nil, testApplicationLoadBalancerFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, nil, nil)
			err := RunApplicationLoadBalancerFlowLogGet(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, nil, testApplicationLoadBalancerFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogGet(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateFlowLog(
				testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, nil, nil)
			err := RunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogCreateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateFlowLog(
				testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, &testResponse, nil)
			err := RunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateFlowLog(
				testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, nil, testApplicationLoadBalancerFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateFlowLog(
				testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testFlowLog, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunApplicationLoadBalancerFlowLogCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties, testQueryParamOther,
			).Return(&testFlowLogUpdated, nil, nil)
			err := RunApplicationLoadBalancerFlowLogUpdate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties, testQueryParamOther,
			).Return(&testFlowLogUpdated, nil, testApplicationLoadBalancerFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogUpdate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogNewVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, &testFlowLogProperties, testQueryParamOther,
			).Return(&testFlowLogUpdated, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunApplicationLoadBalancerFlowLogUpdate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(nil, nil)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, &testResponse, nil)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, &testResponse, testFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListFlowLogs(
				testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(testFlowLogs, &testResponse, nil)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, testFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(nil, testApplicationLoadBalancerFlowLogErr)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgForce, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			cfg.Stdin = bytes.NewReader([]byte("YES\n"))
			rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteFlowLog(
				testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(nil, nil)
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunApplicationLoadBalancerFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testFlowLogVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
			cfg.Stdin = os.Stdin
			err := RunApplicationLoadBalancerFlowLogDelete(cfg)
			assert.Error(t, err)
		},
	)
}
