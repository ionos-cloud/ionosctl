package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testAPIVersion = resources.APIVersion{
		APIVersion: psql.APIVersion{
			Name:       &testAPIVersionVar,
			SwaggerUrl: &testSwaaggerUrlVar,
		},
	}
	testAPIVersions = resources.APIVersionList{
		Versions: []psql.APIVersion{testAPIVersion.APIVersion},
	}
	testAPIVersionVar  = "test-api-version"
	testSwaaggerUrlVar = "/postgresql/test/test/test"
	testAPIVersionErr  = errors.New("test api-version error")
)

func TestAPIVersionCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(APIVersionCmd())
	if ok := APIVersionCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestRunAPIVersionList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), defaultAPIVersionCols)
		rm.CloudApiDbaasPgsqlMocks.Info.EXPECT().List().Return(testAPIVersions, nil, nil)
		err := RunAPIVersionList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAPIVersionListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		rm.CloudApiDbaasPgsqlMocks.Info.EXPECT().List().Return(testAPIVersions, nil, testAPIVersionErr)
		err := RunAPIVersionList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAPIVersionGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		rm.CloudApiDbaasPgsqlMocks.Info.EXPECT().Get().Return(testAPIVersion, nil, nil)
		err := RunAPIVersionGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAPIVersionGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		rm.CloudApiDbaasPgsqlMocks.Info.EXPECT().Get().Return(testAPIVersion, nil, testAPIVersionErr)
		err := RunAPIVersionGet(cfg)
		assert.Error(t, err)
	})
}
