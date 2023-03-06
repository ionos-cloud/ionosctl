package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		err := PreRunDataCenterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDataCenterId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterListFilter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunDataCenterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, nil)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, nil, testDatacenterErr)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&dcNew, &testResponse, testDatacenterErr)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}
