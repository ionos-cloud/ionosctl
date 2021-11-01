package pg

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/dbaas-pg"
	"github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testVersions = resources.PostgresVersionList{
		PostgresVersionList: sdkgo.PostgresVersionList{
			Data: &[]sdkgo.PostgresVersionListData{{
				Name: &testVersionVar,
			}},
		},
	}
	testVersionVar = "test-version"
	testVersionErr = errors.New("test version error")
)

func TestPgsqlVersionCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(PgsqlVersionCmd())
	if ok := PgsqlVersionCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestRunPgsqlVersionGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultPgsqlVersionCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testVersionVar)
		rm.CloudApiDbaasPgsqlMocks.Version.EXPECT().Get(testVersionVar).Return(testVersions, nil, nil)
		err := RunPgsqlVersionGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPgsqlVersionGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultPgsqlVersionCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testVersionVar)
		rm.CloudApiDbaasPgsqlMocks.Version.EXPECT().Get(testVersionVar).Return(testVersions, nil, testVersionErr)
		err := RunPgsqlVersionGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunPgsqlVersionList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultPgsqlVersionCols)
		rm.CloudApiDbaasPgsqlMocks.Version.EXPECT().List().Return(testVersions, nil, nil)
		err := RunPgsqlVersionList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPgsqlVersionListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapidbaaspgsql.ArgClusterId), testVersionVar)
		rm.CloudApiDbaasPgsqlMocks.Version.EXPECT().List().Return(testVersions, nil, testVersionErr)
		err := RunPgsqlVersionList(cfg)
		assert.Error(t, err)
	})
}

func TestGetPgsqlVersionCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("version", config.ArgCols), []string{"Name"})
	getPgsqlVersionCols(core.GetGlobalFlagName("version", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetPgsqlVersionColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("version", config.ArgCols), []string{"Unknown"})
	getPgsqlVersionCols(core.GetGlobalFlagName("version", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
