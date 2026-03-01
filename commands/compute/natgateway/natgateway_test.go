package natgateway

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testutil.TestStateVar},
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
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		err := PreRunNATGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNATGatewayListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testutil.TestQueryParamVar))
		err := PreRunNATGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNATGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		err := PreRunNATGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdsNatGatewayProperties(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayVar)
		err := PreRunDcIdsNatGatewayIps(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdsNatGatewayPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcIdsNatGatewayIps(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		err := PreRunDcNatGatewayIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcNatGatewayIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List().Return(testutil.TestDcs, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testutil.TestDatacenterVar).Return(natgatewaysList, &testutil.TestResponse, nil).Times(len(helpers.GetDataCenters(testutil.TestDcs)))
		err := RunNatGatewayListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.Resource, constants.ArgCols), defaultNatGatewayCols)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgateways, &testutil.TestResponse, nil)
		err := RunNatGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.Resource, constants.ArgCols), defaultNatGatewayCols)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagOrderBy), testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(resources.NatGateways{}, &testutil.TestResponse, nil)
		err := RunNatGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgateways, nil, testNatGatewayErr)
		err := RunNatGatewayList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, &testutil.TestResponse, nil)
		err := RunNatGatewayGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, nil)
		err := RunNatGatewayGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, testNatGatewayErr)
		err := RunNatGatewayGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, testNatGatewayErr)
		err := RunNatGatewayGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testutil.TestResponse, nil)
		err := RunNatGatewayCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testutil.TestResponse, testNatGatewayErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, nil, testNatGatewayErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, &testutil.TestResponse, nil)
		err := RunNatGatewayUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayNewVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, nil, testNatGatewayErr)
		err := RunNatGatewayUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayNewVar})
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgatewaysList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgatewaysList, nil, testNatGatewayErr)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(resources.NatGateways{}, &testutil.TestResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(
			resources.NatGateways{NatGateways: ionoscloud.NatGateways{Items: &[]ionoscloud.NatGateway{}}}, &testutil.TestResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgatewaysList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, testNatGatewayErr)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, testNatGatewayErr)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}
