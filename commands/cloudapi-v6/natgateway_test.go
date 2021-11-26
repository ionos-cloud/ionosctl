package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
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
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunNATGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNATGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunNATGatewayList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcIdsNatGatewayProperties(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcIdsNatGatewayIps(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNatGatewayIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgCols), defaultNatGatewayCols)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar, resources.ListQueryParams{}).Return(natgateways, &testResponse, nil)
		err := RunNatGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgCols), defaultNatGatewayCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar, testListQueryParam).Return(resources.NatGateways{}, &testResponse, nil)
		err := RunNatGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar, resources.ListQueryParams{}).Return(natgateways, nil, testNatGatewayErr)
		err := RunNatGatewayList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, &testResponse, nil)
		err := RunNatGatewayGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testResponse, nil)
		err := RunNatGatewayCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testResponse, testNatGatewayErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayNewVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, &testResponse, nil)
		err := RunNatGatewayUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().List(testNatGatewayVar, resources.ListQueryParams{}).Return(natgatewaysList, &testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testResponse, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayVar)
		cfg.Stdin = os.Stdin
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNatGatewaysCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Name"})
	getNatGatewaysCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNatGatewaysColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Unknown"})
	getNatGatewaysCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
