package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testFlowLog = resources.FlowLog{
		FlowLog: compute.FlowLog{
			Id: &testFlowLogVar,
			Properties: &compute.FlowLogProperties{
				Name:      &testFlowLogVar,
				Action:    &testFlowLogUpperVar,
				Direction: &testFlowLogUpperVar,
				Bucket:    &testFlowLogVar,
			},
			Metadata: &compute.DatacenterElementMetadata{
				State: &testFlowLogState,
			},
		},
	}
	testFlowLogsList = resources.FlowLogs{
		FlowLogs: compute.FlowLogs{
			Id: &testFlowLogVar,
			Items: &[]compute.FlowLog{
				testFlowLog.FlowLog,
				testFlowLog.FlowLog,
			},
		},
	}
	testInputFlowLog = resources.FlowLog{
		FlowLog: compute.FlowLog{
			Properties: testFlowLog.FlowLog.Properties,
		},
	}
	testFlowLogUpdated = resources.FlowLog{
		FlowLog: compute.FlowLog{
			Properties: &testFlowLogProperties.FlowLogProperties,
		},
	}
	testFlowLogProperties = resources.FlowLogProperties{
		FlowLogProperties: compute.FlowLogProperties{
			Name:      &testFlowLogNewVar,
			Action:    &testFlowLogNewUpperVar,
			Direction: &testFlowLogNewUpperVar,
			Bucket:    &testFlowLogNewVar,
		},
	}
	testFlowLogs = resources.FlowLogs{
		FlowLogs: compute.FlowLogs{
			Id:    &testFlowLogVar,
			Items: &[]compute.FlowLog{testFlowLog.FlowLog},
		},
	}
	testFlowLogState       = "AVAILABLE"
	testFlowLogVar         = "test-flowlog"
	testFlowLogUpperVar    = strings.ToUpper(testFlowLogVar)
	testFlowLogNewVar      = "test-new-flowlog"
	testFlowLogNewUpperVar = strings.ToUpper(testFlowLogNewVar)
	testFlowLogErr         = errors.New("flowlog test error")
)

func TestFlowlogCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(FlowlogCmd())
	if ok := FlowlogCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		err := PreRunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFlowLogListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFlowLogListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		err := PreRunFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerNicFlowLogIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		err := PreRunDcServerNicFlowLogIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerNicFlowLogIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcServerNicFlowLogIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFlowLogs, &testResponse, nil)
		err := RunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.FlowLogs{}, &testResponse, nil)
		err := RunFlowLogList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFlowLogs, nil, testFlowLogErr)
		err := RunFlowLogList(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Get(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testFlowLog, nil, testFlowLogErr)
		err := RunFlowLogGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Get(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testFlowLog, &testResponse, nil)
		err := RunFlowLogGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Create(testFlowLogVar, testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testInputFlowLog, &testResponse, nil)
		err := RunFlowLogCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAction), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Bucket), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Create(testFlowLogVar, testFlowLogVar, testFlowLogVar, testInputFlowLog, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testInputFlowLog, &testResponse, testFlowLogErr)
		err := RunFlowLogCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFlowLogsList, &testResponse, nil)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFlowLogsList, &testResponse, testFlowLogErr)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.FlowLogs{}, &testResponse, nil)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.FlowLogs{FlowLogs: compute.FlowLogs{Items: &[]compute.FlowLog{}}}, &testResponse, nil)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().List(testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFlowLogsList, &testResponse, nil)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testFlowLogErr)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testFlowLogErr)
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFlowLogDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FlowLog.EXPECT().Delete(testFlowLogVar, testFlowLogVar, testFlowLogVar, testFlowLogVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunFlowLogDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFlowLogDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFlowLogId), testFlowLogVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunFlowLogDelete(cfg)
		assert.Error(t, err)
	})
}
