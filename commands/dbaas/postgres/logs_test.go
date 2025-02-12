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
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLogs = resources.ClusterLogs{
		ClusterLogs: psql.ClusterLogs{
			Instances: &[]psql.ClusterLogsInstancesInner{{
				Name: &testLogVar,
				Messages: &[]psql.ClusterLogsInstancesInnerMessagesInner{
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSince), testSinceVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgUntil), testUntilVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSince), "3min")
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgUntil), testUntilVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSince), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgUntil), "1min")
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDirection), testDirectionVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgSince), testSinceVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgUntil), testUntilVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDirection), testDirectionVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLimit), testLimitVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgDirection), testDirectionVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, &testLogsQueryParams).Return(&testLogs, nil, testLogErr)
		err := RunClusterLogsList(cfg)
		assert.Error(t, err)
	})
}
