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
	natgatewayTest = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Properties: &ionoscloud.NatGatewayProperties{
				Name:      &testNatGatewayVar,
				PublicIps: &[]string{testNatGatewayVar},
			},
		},
	}
	natgatewayTestId = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Id: &testNatGatewayVar,
			Properties: &ionoscloud.NatGatewayProperties{
				Name:      &testNatGatewayVar,
				PublicIps: &[]string{testNatGatewayVar},
			},
		},
	}
	natgatewayTestGet = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Id:         &testNatGatewayVar,
			Properties: natgatewayTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgateways = resources.NatGateways{
		NatGateways: ionoscloud.NatGateways{
			Id:    &testNatGatewayVar,
			Items: &[]ionoscloud.NatGateway{natgatewayTest.NatGateway},
		},
	}
	natgatewaysList = resources.NatGateways{
		NatGateways: ionoscloud.NatGateways{
			Id: &testNatGatewayVar,
			Items: &[]ionoscloud.NatGateway{
				natgatewayTestId.NatGateway,
				natgatewayTestId.NatGateway,
			},
		},
	}
	natgatewayProperties = resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Name:      &testNatGatewayNewVar,
			PublicIps: &[]string{testNatGatewayNewVar},
		},
	}
	natgatewayNew = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Properties: &natgatewayProperties.NatGatewayProperties,
		},
	}
	testNatGatewayVar    = "test-natgateway"
	testNatGatewayNewVar = "test-new-natgateway"
	testNatGatewayErr    = errors.New("natgateway test error")
)

func TestNatgatewayCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(NatgatewayCmd())
	if ok := NatgatewayCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunNATGatewayList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			err := PreRunNATGatewayList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunNATGatewayListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(
				core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters),
				[]string{fmt.Sprintf("createdBy=%s", testQueryParamVar)},
			)
			err := PreRunNATGatewayList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunNATGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(
				core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters),
				[]string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)},
			)
			err := PreRunNATGatewayList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestPreRunDcIdsNatGatewayProperties(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayVar)
			err := PreRunDcIdsNatGatewayIps(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunDcIdsNatGatewayPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			err := PreRunDcIdsNatGatewayIps(cfg)
			assert.Error(t, err)
		},
	)
}

func TestPreRunDcNatGatewayIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			err := PreRunDcNatGatewayIds(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunDcNatGatewayIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			err := PreRunDcNatGatewayIds(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
				dcs, &testResponse, nil,
			)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testDatacenterVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgatewaysList, &testResponse, nil).Times(len(getDataCenters(dcs)))
			err := RunNatGatewayListAll(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.Resource, constants.ArgCols), defaultNatGatewayCols)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgateways, &testResponse, nil)
			err := RunNatGatewayList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.Resource, constants.ArgCols), defaultNatGatewayCols)
			viper.Set(
				core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters),
				[]string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)},
			)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(resources.NatGateways{}, &testResponse, nil)
			err := RunNatGatewayList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgateways, nil, testNatGatewayErr)
			err := RunNatGatewayList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTestGet, &testResponse, nil)
			err := RunNatGatewayGet(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTestGet, nil, nil)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTestGet, nil, nil)
			err := RunNatGatewayGet(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTestGet, nil, testNatGatewayErr)
			err := RunNatGatewayGet(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTestGet, nil, testNatGatewayErr)
			err := RunNatGatewayGet(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(
				testNatGatewayVar, natgatewayTest, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTest, &testResponse, nil)
			err := RunNatGatewayCreate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(
				testNatGatewayVar, natgatewayTest, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTest, &testResponse, testNatGatewayErr)
			err := RunNatGatewayCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(
				testNatGatewayVar, natgatewayTest, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTest, nil, testNatGatewayErr)
			err := RunNatGatewayCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(
				testNatGatewayVar, natgatewayTest, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayTest, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunNatGatewayCreate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayNewVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(
				testNatGatewayVar, testNatGatewayVar, natgatewayProperties,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayNew, &testResponse, nil)
			err := RunNatGatewayUpdate(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayNewVar})
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(
				testNatGatewayVar, testNatGatewayVar, natgatewayProperties,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayNew, nil, testNatGatewayErr)
			err := RunNatGatewayUpdate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayNewVar})
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(
				testNatGatewayVar, testNatGatewayVar, natgatewayProperties,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&natgatewayNew, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunNatGatewayUpdate(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			err := RunNatGatewayDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgatewaysList, &testResponse, nil)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			err := RunNatGatewayDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgatewaysList, nil, testNatGatewayErr)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(resources.NatGateways{}, &testResponse, nil)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(
				resources.NatGateways{NatGateways: ionoscloud.NatGateways{Items: &[]ionoscloud.NatGateway{}}},
				&testResponse, nil,
			)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().List(
				testNatGatewayVar, gomock.AssignableToTypeOf(testListQueryParam),
			).Return(natgatewaysList, &testResponse, nil)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, testNatGatewayErr)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(nil, testNatGatewayErr)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgForce, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			cfg.Stdin = bytes.NewReader([]byte("YES\n"))
			rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(
				testNatGatewayVar, testNatGatewayVar, gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(nil, nil)
			err := RunNatGatewayDelete(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNatGatewayDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgForce, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
			cfg.Stdin = os.Stdin
			err := RunNatGatewayDelete(cfg)
			assert.Error(t, err)
		},
	)
}
