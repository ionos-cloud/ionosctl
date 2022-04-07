package postgres

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testAPIVersion = resources.APIVersion{
		APIVersion: sdkgo.APIVersion{
			Name:       &testAPIVersionVar,
			SwaggerUrl: &testSwaaggerUrlVar,
		},
	}
	testAPIVersions = resources.APIVersionList{
		Versions: []sdkgo.APIVersion{testAPIVersion.APIVersion},
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.NS, config.ArgCols), defaultAPIVersionCols)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.CloudApiDbaasPgsqlMocks.Info.EXPECT().Get().Return(testAPIVersion, nil, testAPIVersionErr)
		err := RunAPIVersionGet(cfg)
		assert.Error(t, err)
	})
}

func TestGetAPIVersionsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("api-version", config.ArgCols), []string{"SwaggerUrl"})
	getAPIVersionCols(core.GetGlobalFlagName("api-version", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetAPIVersionsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("api-version", config.ArgCols), []string{"Unknown"})
	getAPIVersionCols(core.GetGlobalFlagName("api-version", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
