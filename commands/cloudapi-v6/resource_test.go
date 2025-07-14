package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	resourceTest = resources.Resource{
		Resource: compute.Resource{
			Properties: &compute.ResourceProperties{
				Name:              &testResourceVar,
				SecAuthProtection: &testResourceBoolVar,
			},
			Type: &testResourceType,
		},
	}
	resourceTestGet = resources.Resource{
		Resource: compute.Resource{
			Id: &testResourceVar,
			Properties: &compute.ResourceProperties{
				Name:              &testResourceVar,
				SecAuthProtection: &testResourceBoolVar,
			},
			Metadata: &compute.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	rs = resources.Resources{
		Resources: compute.Resources{
			Id:    &testResourceVar,
			Items: &[]compute.Resource{resourceTest.Resource},
		},
	}
	resourceGroupTest = resources.ResourceGroups{
		ResourceGroups: compute.ResourceGroups{
			Id:    &testResourceVar,
			Items: &[]compute.Resource{resourceTest.Resource},
		},
	}
	testResourceType    = compute.Type(testResourceVar)
	testResourceBoolVar = false
	testResourceVar     = "test-resource"
	testResourceErr     = errors.New("resource test error")
)

func TestResourceCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ResourceCmd())
	if ok := ResourceCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunResourceType(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testResourceVar)
		err := PreRunResourceType(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunResourceType(cfg)
		assert.Error(t, err)
	})
}

func TestRunResourceList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.User.EXPECT().ListResources().Return(rs, &testResponse, nil)
		err := RunResourceList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourcesByType(testResourceVar).Return(rs, &testResponse, nil)
		err := RunResourceGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceGetByTypeAndId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testResourceVar)
		rm.CloudApiV6Mocks.User.EXPECT().GetResourceByTypeAndId(testResourceVar, testResourceVar).Return(&resourceTestGet, &testResponse, nil)
		err := RunResourceGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunResourceGetByTypeAndIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testResourceVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testResourceVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testResourceVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListResources(testResourceVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resourceGroupTest, &testResponse, nil)
		err := RunGroupResourceList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupResourceListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testResourceVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListResources(testResourceVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resourceGroupTest, nil, testResourceErr)
		err := RunGroupResourceList(cfg)
		assert.Error(t, err)
	})
}
