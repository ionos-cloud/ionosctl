package cloudapi_dbaas_pgsql

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
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
						Time:    &testLogVar,
						Message: &testLogVar,
					},
				},
			}},
		},
	}
	testLimitVar = int32(1)
	testLogVar   = "test-cluster-logs"
	testLogErr   = errors.New("test cluster-logs error")
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStartTime), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgEndTime), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLimit), testLimitVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, testLogVar, testLogVar, testLimitVar).Return(testLogs, nil, nil)
		err := RunClusterLogsGet(cfg)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgStartTime), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgEndTime), testLogVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgLimit), testLimitVar)
		rm.CloudApiDbaasPgsqlMocks.Log.EXPECT().Get(testLogVar, testLogVar, testLogVar, testLimitVar).Return(testLogs, nil, testLogErr)
		err := RunClusterLogsGet(cfg)
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
