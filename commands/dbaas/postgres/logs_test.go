package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLogs = resources.ClusterLogs{
		ClusterLogs: sdkgo.ClusterLogs{
			Instances: &[]sdkgo.ClusterLogsInstancesInner{{
				Name: &testLogVar,
				Messages: &[]sdkgo.ClusterLogsInstancesInnerMessagesInner{
					{
						Time:    &testIonosTime,
						Message: &testLogVar,
					},
				},
			}},
		},
	}
	testStartTimeVar    = "2021-01-01T00:00:00Z"
	testEndTimeVar      = "2021-02-02T00:00:00Z"
	testSinceVar        = "2h"
	testUntilVar        = "1h"
	testStartTime       = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	testEndTime         = time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)
	testLimitVar        = int32(1)
	testDirectionVar    = "BACKWARD"
	testLogVar          = "test-cluster-logs"
	testLogErr          = errors.New("test cluster-logs error")
	testLogsQueryParams = resources.LogsQueryParams{
		Direction: testDirectionVar,
		Limit:     testLimitVar,
		StartTime: testStartTime,
		EndTime:   testEndTime,
	}
)

func TestLogCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LogsCmd())
	if ok := LogsCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunClusterLogsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagSince), testSinceVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagUntil), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		err := PreRunClusterLogsList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunClusterLogsListSinceErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagSince), "3min")
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagUntil), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		err := PreRunClusterLogsList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunClusterLogsListUntilErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagSince), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagUntil), "1min")
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testClusterVar)
		err := PreRunClusterLogsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunClusterLogsGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagDirection), testDirectionVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, &testLogsQueryParams).Return(&testLogs, nil, nil)
		err := RunClusterLogsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterLogsGetSinceUntilIgnored(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagSince), testSinceVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagUntil), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagDirection), testDirectionVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, &testLogsQueryParams).Return(&testLogs, nil, nil)
		err := RunClusterLogsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterLogsGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.FlagDirection), testDirectionVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, &testLogsQueryParams).Return(&testLogs, nil, testLogErr)
		err := RunClusterLogsList(cfg)
		assert.Error(t, err)
	})
}
