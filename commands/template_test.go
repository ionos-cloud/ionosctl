package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	tpl = v6.Template{
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
	templates = v6.Templates{
		Templates: ionoscloud.Templates{
			Items: &[]ionoscloud.Template{tpl.Template},
		},
	}
	testTemplateSize = float32(2)
	testTemplateVar  = "test-template"
	testTemplateErr  = errors.New("template test error")
)

func TestPreTemplateId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTemplateId), testTemplateVar)
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.Template.EXPECT().List().Return(templates, &testResponse, nil)
		err := RunTemplateList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTemplateListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Template.EXPECT().List().Return(templates, nil, testTemplateErr)
		err := RunTemplateList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTemplateGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTemplateId), testTemplateVar)
		rm.Template.EXPECT().Get(testTemplateVar).Return(&tpl, &testResponse, nil)
		err := RunTemplateGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTemplateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTemplateId), testTemplateVar)
		rm.Template.EXPECT().Get(testTemplateVar).Return(&tpl, nil, testTemplateErr)
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

func TestGetTemplatesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getTemplatesIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
