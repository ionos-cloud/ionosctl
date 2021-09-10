package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	resourceTest = resources.Resource{
		Resource: ionoscloud.Resource{
			Properties: &ionoscloud.ResourceProperties{
				Name:              &testResourceVar,
				SecAuthProtection: &testResourceBoolVar,
			},
			Type: &testResourceType,
		},
	}
	resourceTestGet = resources.Resource{
		Resource: ionoscloud.Resource{
			Id: &testResourceVar,
			Properties: &ionoscloud.ResourceProperties{
				Name:              &testResourceVar,
				SecAuthProtection: &testResourceBoolVar,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	rs = resources.Resources{
		Resources: ionoscloud.Resources{
			Id:    &testResourceVar,
			Items: &[]ionoscloud.Resource{resourceTest.Resource},
		},
	}
	resourceGroupTest = resources.ResourceGroups{
		ResourceGroups: ionoscloud.ResourceGroups{
			Id:    &testResourceVar,
			Items: &[]ionoscloud.Resource{resourceTest.Resource},
		},
	}
	testResourceType    = ionoscloud.Type(testResourceVar)
	testResourceBoolVar = false
	testResourceVar     = "test-resource"
	testResourceErr     = errors.New("resource test error")
)

func TestPreRunResourceType(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgType), testResourceVar)
		err := PreRunResourceType(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunResourceType(cfg)
		assert.Error(t, err)
	})
}

func TestRunResourceList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.User.EXPECT().ListResources().Return(rs, nil, nil)
		err := RunResourceList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.User.EXPECT().ListResources().Return(rs, nil, testResourceErr)
		err := RunResourceList(cfg)
		assert.Error(t, err)
	})
}

func TestRunResourceGetByType(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgType), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourcesByType(testResourceVar).Return(rs, nil, nil)
		err := RunResourceGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceGetByTypeAndId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgType), testResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgResourceId), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourceByTypeAndId(testResourceVar, testResourceVar).Return(&resourceTestGet, nil, nil)
		err := RunResourceGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceGetByTypeAndIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgType), testResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgResourceId), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourceByTypeAndId(testResourceVar, testResourceVar).Return(&resourceTestGet, nil, testResourceErr)
		err := RunResourceGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunResourceGetByTypeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgType), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourcesByType(testResourceVar).Return(rs, nil, testResourceErr)
		err := RunResourceGet(cfg)
		assert.Error(t, err)
	})
}

// Group Resources

func TestRunGroupResourceList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgGroupId), testResourceVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListResources(testResourceVar).Return(resourceGroupTest, nil, nil)
		err := RunGroupResourceList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupResourceListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgGroupId), testResourceVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListResources(testResourceVar).Return(resourceGroupTest, nil, testResourceErr)
		err := RunGroupResourceList(cfg)
		assert.Error(t, err)
	})
}

func TestGetResourcesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("resource", config.ArgCols), []string{"Type"})
	getResourceCols(core.GetGlobalFlagName("resource", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetResourcesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("resource", config.ArgCols), []string{"Unknown"})
	getResourceCols(core.GetGlobalFlagName("resource", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetResourcesIds(t *testing.T) {
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
	getResourcesIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
