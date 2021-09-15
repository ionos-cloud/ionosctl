package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
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
		rm.CloudApiV6Mocks.Template.EXPECT().List().Return(templates, nil, nil)
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
		rm.CloudApiV6Mocks.Template.EXPECT().List().Return(templates, nil, testTemplateErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testTemplateVar)
		rm.CloudApiV6Mocks.Template.EXPECT().Get(testTemplateVar).Return(&tpl, nil, nil)
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
		rm.CloudApiV6Mocks.Template.EXPECT().Get(testTemplateVar).Return(&tpl, nil, testTemplateErr)
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
