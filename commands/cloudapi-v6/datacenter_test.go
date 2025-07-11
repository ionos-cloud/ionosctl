package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	dcVersion = int32(1)
	dc        = ionoscloud.Datacenter{
		Id: &testDatacenterVar,
		Properties: &ionoscloud.DatacenterProperties{
			Name:        &testDatacenterVar,
			Description: &testDatacenterVar,
			Location:    &testDatacenterVar,
			Version:     &dcVersion,
			Features:    &[]string{testDatacenterVar},
			CpuArchitecture: &[]ionoscloud.CpuArchitectureProperties{
				{
					CpuFamily: &testDatacenterVar,
				}},
			SecAuthProtection: &testDatacenterBoolVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testStateVar,
		},
	}
	dcProperties = resources.DatacenterProperties{
		DatacenterProperties: ionoscloud.DatacenterProperties{
			Name:        &testDatacenterNewVar,
			Description: &testDatacenterNewVar,
		},
	}
	dcPropertiesPut = resources.DatacenterPropertiesPut{
		DatacenterPropertiesPut: ionoscloud.DatacenterPropertiesPut{
			Name:        &testDatacenterNewVar,
			Description: &testDatacenterNewVar,
		},
	}
	dcNew = resources.Datacenter{
		Datacenter: ionoscloud.Datacenter{
			Id: &testDatacenterVar,
			Properties: &ionoscloud.DatacenterProperties{
				Name:        dcProperties.DatacenterProperties.Name,
				Description: dcProperties.DatacenterProperties.Description,
				Location:    &testDatacenterVar,
			},
		},
	}
	dcs = resources.Datacenters{
		Datacenters: ionoscloud.Datacenters{
			Id:    &testDatacenterVar,
			Items: &[]ionoscloud.Datacenter{dc, dc},
		},
	}
	testDatacenterVar     = "test-datacenter"
	testDatacenterBoolVar = false
	testDatacenterNewVar  = "test-new-datacenter"
	testDatacenterErr     = errors.New("datacenter test error occurred")
)

func TestDatacenterCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DatacenterCmd())
	if ok := DatacenterCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDataCenterId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		err := PreRunDataCenterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunDataCenterId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterListFilter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, &testResponse, nil)
		err := RunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Datacenters{}, &testResponse, nil)
		err := RunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, nil, testDatacenterErr)
		err := RunDataCenterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		res := resources.Datacenter{Datacenter: dc}
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&res, &testResponse, nil)
		err := RunDataCenterGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		res := resources.Datacenter{Datacenter: dc}
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Get(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&res, nil, testDatacenterErr)
		err := RunDataCenterGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Datacenter{dc}, &testResponse, nil)
		err := RunDataCenterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Datacenter{Datacenter: dc}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunDataCenterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Datacenter{Datacenter: dc}, nil, testDatacenterErr)
		err := RunDataCenterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcPropertiesPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcPropertiesPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcPropertiesPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcPropertiesPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, nil, testDatacenterErr)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Datacenters{}, nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Datacenters{Datacenters: ionoscloud.Datacenters{Items: &[]ionoscloud.Datacenter{}}}, &testResponse, nil)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testDatacenterErr)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testDatacenterErr)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Delete(testDatacenterVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}
