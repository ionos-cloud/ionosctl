package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	tpl = resources.Template{
		Template: ionoscloud.Template{
			Id: &testTemplateVar,
			Properties: &ionoscloud.TemplateProperties{
				Name:        &testTemplateVar,
				Cores:       &testTemplateSize,
				Ram:         &testTemplateSize,
				StorageSize: &testTemplateSize,
			},
		},
	}
	templates = resources.Templates{
		Templates: ionoscloud.Templates{
			Items: &[]ionoscloud.Template{tpl.Template},
		},
	}
	testTemplateSize = float32(2)
	testTemplateVar  = "test-template"
	testTemplateErr  = errors.New("template test error")
)

func TestTemplateCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(TemplateCmd())
	if ok := TemplateCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunTemplateList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunTemplateList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTemplateListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunTemplateList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTemplateListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunTemplateList(cfg)
		assert.Error(t, err)
	})
}

func TestPreTemplateId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testTemplateVar)
		err := PreRunTemplateId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreTemplateIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunTemplateId(cfg)
		assert.Error(t, err)
	})
}

func TestRunTemplateList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		rm.CloudApiV6Mocks.Template.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(templates, &testResponse, nil)
		err := RunTemplateList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTemplateListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Template.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Templates{}, &testResponse, nil)
		err := RunTemplateList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTemplateListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.Template.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(templates, nil, testTemplateErr)
		err := RunTemplateList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTemplateGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testTemplateVar)
		rm.CloudApiV6Mocks.Template.EXPECT().Get(testTemplateVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&tpl, &testResponse, nil)
		err := RunTemplateGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTemplateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testTemplateVar)
		rm.CloudApiV6Mocks.Template.EXPECT().Get(testTemplateVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&tpl, nil, testTemplateErr)
		err := RunTemplateGet(cfg)
		assert.Error(t, err)
	})
}

func TestGetTemplatesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("template", config.ArgCols), []string{"Name"})
	getTemplateCols(core.GetGlobalFlagName("template", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetTemplatesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("template", config.ArgCols), []string{"Unknown"})
	getTemplateCols(core.GetGlobalFlagName("template", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
