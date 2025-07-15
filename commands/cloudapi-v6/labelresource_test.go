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
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLabelResource = compute.LabelResource{
		Id: &testLabelVar,
		Properties: compute.LabelResourceProperties{
			Key:   &testLabelResourceVar,
			Value: &testLabelResourceVar,
		},
	}
	testLabelResources = resources.LabelResources{
		LabelResources: compute.LabelResources{
			Id:    &testLabelVar,
			Items: &[]compute.LabelResource{testLabelResource},
		},
	}
	testLabelResourcesList = resources.LabelResources{
		LabelResources: compute.LabelResources{
			Id: &testLabelVar,
			Items: &[]compute.LabelResource{
				testLabelResource,
				testLabelResource,
			},
		},
	}
	testLabelResourceRes = resources.LabelResource{LabelResource: testLabelResource}
	testLabelResourceVar = "test-label-resource"
	testLabelResourceErr = errors.New("label resource test error")
)

func TestRunDataCenterLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunDataCenterLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunDataCenterLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunDataCenterLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunDataCenterLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, nil, testLabelResourceErr)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(resources.LabelResources{}, &testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(
			resources.LabelResources{LabelResources: compute.LabelResources{Items: &[]compute.LabelResource{}}}, &testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, testLabelResourceErr)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunIpBlockLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunIpBlockLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunIpBlockLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunIpBlockLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, nil, testLabelResourceErr)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(resources.LabelResources{}, &testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(
			resources.LabelResources{LabelResources: compute.LabelResources{Items: &[]compute.LabelResource{}}}, &testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, testLabelResourceErr)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunSnapshotLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunSnapshotLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunSnapshotLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunSnapshotLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, nil, testLabelResourceErr)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(resources.LabelResources{}, &testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(
			resources.LabelResources{LabelResources: compute.LabelResources{Items: &[]compute.LabelResource{}}}, &testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, testLabelResourceErr)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunServerLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunServerLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunServerLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).
			Return(&testLabelResourceRes, &testResponse, nil)
		err := RunServerLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		// viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(resources.LabelResources{}, &testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, nil, testLabelResourceErr)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(
			resources.LabelResources{LabelResources: compute.LabelResources{Items: &[]compute.LabelResource{}}}, &testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, testLabelResourceErr)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunVolumeLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunVolumeLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunVolumeLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).
			Return(&testLabelResourceRes, &testResponse, nil)
		err := RunVolumeLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, nil, testLabelResourceErr)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(resources.LabelResources{}, &testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(
			resources.LabelResources{LabelResources: compute.LabelResources{Items: &[]compute.LabelResource{}}}, &testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResourcesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, testLabelResourceErr)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}
