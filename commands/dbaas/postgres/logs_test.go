package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLogs = resources.ClusterLogs{
		ClusterLogs: sdkgo.ClusterLogs{
			Instances: &[]sdkgo.ClusterLogsInstances{{
				Name: &testLogVar,
				Messages: &[]sdkgo.ClusterLogsMessages{
					{
						Time:    &testIonosTime,
						Message: &testLogVar,
					},
				},
			}},
		},
	}
	testStartTimeVar = "2021-01-01T00:00:00Z"
	testEndTimeVar   = "2021-02-02T00:00:00Z"
	testStartTime    = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	testEndTime      = time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)
	testLimitVar     = int32(1)
	testLogVar       = "test-cluster-logs"
	testLogErr       = errors.New("test cluster-logs error")
)

func TestLogCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LogsCmd())
	if ok := LogsCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestRunClusterLogsGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLimit), testLimitVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, testLimitVar, testStartTime, testEndTime).Return(&testLogs, nil, nil)
		err := RunClusterLogsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunClusterLogsGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgStartTime), testStartTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgEndTime), testEndTimeVar)
		viper.Set(core.GetFlagName(cfg.NS, dbaaspg.ArgLimit), testLimitVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, testLimitVar, testStartTime, testEndTime).Return(&testLogs, nil, testLogErr)
		err := RunClusterLogsList(cfg)
		assert.Error(t, err)
	})
}

func TestGetClusterLogsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("logs", config.ArgCols), []string{"Name"})
	getClusterLogsCols(core.GetGlobalFlagName("logs", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLogsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("logs", config.ArgCols), []string{"Unknown"})
	getClusterLogsCols(core.GetGlobalFlagName("logs", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
