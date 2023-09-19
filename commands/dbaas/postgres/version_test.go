package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultPgsqlVersionCols)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testVersionVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultPgsqlVersionCols)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testVersionVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultPgsqlVersionCols)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testVersionVar)
		rm.CloudApiDbaasPgsqlMocks.Version.EXPECT().List().Return(testVersions, nil, testVersionErr)
		err := RunPgsqlVersionList(cfg)
		assert.Error(t, err)
	})
}
