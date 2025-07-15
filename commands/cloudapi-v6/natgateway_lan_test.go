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
	natgatewayLanTest = resources.NatGateway{
		NatGateway: compute.NatGateway{
			Properties: &compute.NatGatewayProperties{
				Name:      &testNatGatewayLanVar,
				PublicIps: &[]string{testNatGatewayLanVar},
				Lans:      &[]compute.NatGatewayLanProperties{natgatewayLanProperties.NatGatewayLanProperties},
			},
		},
	}
	natgatewaysLanListTest = resources.NatGateways{
		NatGateways: compute.NatGateways{
			Id:    &testNatGatewayLanVar,
			Items: &[]compute.NatGateway{natgatewayLanTest.NatGateway, natgatewayLanTest.NatGateway},
		},
	}
	natgatewayLanTestUpdated = resources.NatGateway{
		NatGateway: compute.NatGateway{
			Id: &testNatGatewayLanVar,
			Properties: &compute.NatGatewayProperties{
				Name:      &testNatGatewayLanVar,
				PublicIps: &[]string{testNatGatewayLanVar},
				Lans:      &[]compute.NatGatewayLanProperties{natgatewayLanProperties.NatGatewayLanProperties},
			},
			Metadata: &compute.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgatewayLanTestProper = resources.NatGatewayProperties{
		NatGatewayProperties: compute.NatGatewayProperties{
			Lans: &[]compute.NatGatewayLanProperties{
				natgatewayLanProperties.NatGatewayLanProperties,
				natgatewayLanNewProperties.NatGatewayLanProperties,
			},
		},
	}
	// Send empty struct to overwrite the existing one
	natgatewayLanTestRemove = resources.NatGatewayProperties{
		NatGatewayProperties: compute.NatGatewayProperties{
			Lans: &[]compute.NatGatewayLanProperties{
				natgatewayLanProperties.NatGatewayLanProperties,
			},
		},
	}
	natgatewayLanTestRemoveAll = resources.NatGatewayProperties{
		NatGatewayProperties: compute.NatGatewayProperties{
			Lans: &[]compute.NatGatewayLanProperties{},
		},
	}
	natgatewayLanProperties = resources.NatGatewayLanProperties{
		NatGatewayLanProperties: compute.NatGatewayLanProperties{
			Id:         &testNatGatewayLanIntVar,
			GatewayIps: &[]string{testNatGatewayLanVar},
		},
	}
	natgatewayLanNewProperties = resources.NatGatewayLanProperties{
		NatGatewayLanProperties: compute.NatGatewayLanProperties{
			Id:         &testNatGatewayLanNewIntVar,
			GatewayIps: &[]string{testNatGatewayLanNewVar},
		},
	}
	testNatGatewayLanIntVar    = int32(1)
	testNatGatewayLanNewIntVar = int32(2)
	testNatGatewayLanVar       = "test-natgateway-lan"
	testNatGatewayLanNewVar    = "test-new-natgateway-lan"
	testNatGatewayLanErr       = errors.New("natgateway-lan test error")
)

func TestPreRunDcNatGatewayLanIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanVar)
		err := PreRunDcNatGatewayLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcNatGatewayLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.Resource, constants.ArgCols), defaultNatGatewayLanCols)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, &testResponse, nil)
		err := RunNatGatewayLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNatGatewayLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTestUpdated, &testResponse, nil)
		err := RunNatGatewayLanAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTestUpdated, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTestUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, &testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, &testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemoveAll, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testNatGatewayLanNewIntVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}
