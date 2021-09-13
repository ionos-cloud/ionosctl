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
	natgatewayRuleTest = resources.NatGatewayRule{
		NatGatewayRule: ionoscloud.NatGatewayRule{
			Properties: &ionoscloud.NatGatewayRuleProperties{
				Name:         &testNatGatewayRuleVar,
				PublicIp:     &testNatGatewayRuleVar,
				Protocol:     &testNatGatewayRuleProtocol,
				SourceSubnet: &testNatGatewayRuleVar,
				TargetSubnet: &testNatGatewayRuleVar,
				TargetPortRange: &ionoscloud.TargetPortRange{
					Start: &testNatGatewayRuleIntVar,
					End:   &testNatGatewayRuleIntVar,
				},
			},
		},
	}
	natgatewayRuleTestGet = resources.NatGatewayRule{
		NatGatewayRule: ionoscloud.NatGatewayRule{
			Id:         &testNatGatewayRuleVar,
			Properties: natgatewayRuleTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgatewayRules = resources.NatGatewayRules{
		NatGatewayRules: ionoscloud.NatGatewayRules{
			Id:    &testNatGatewayRuleVar,
			Items: &[]ionoscloud.NatGatewayRule{natgatewayRuleTestGet.NatGatewayRule},
		},
	}
	natgatewayRuleProperties = resources.NatGatewayRuleProperties{
		NatGatewayRuleProperties: ionoscloud.NatGatewayRuleProperties{
			Name:         &testNatGatewayRuleNewVar,
			PublicIp:     &testNatGatewayRuleNewVar,
			Protocol:     &testNatGatewayRuleNewProtocol,
			SourceSubnet: &testNatGatewayRuleNewVar,
			TargetSubnet: &testNatGatewayRuleNewVar,
			TargetPortRange: &ionoscloud.TargetPortRange{
				Start: &testNatGatewayRuleNewIntVar,
				End:   &testNatGatewayRuleNewIntVar,
			},
		},
	}
	natgatewayRuleNew = resources.NatGatewayRule{
		NatGatewayRule: ionoscloud.NatGatewayRule{
			Id:         &testNatGatewayRuleVar,
			Properties: &natgatewayRuleProperties.NatGatewayRuleProperties,
		},
	}
	testNatGatewayRuleIntVar      = int32(10000)
	testNatGatewayRuleNewIntVar   = int32(20000)
	testNatGatewayRuleProtocol    = ionoscloud.NatGatewayRuleProtocol("ALL")
	testNatGatewayRuleNewProtocol = ionoscloud.NatGatewayRuleProtocol("TCP")
	testNatGatewayRuleVar         = "test-natgateway-rule"
	testNatGatewayRuleNewVar      = "test-new-natgateway-rule"
	testNatGatewayRuleErr         = errors.New("natgateway-rule test error")
)

func TestPreRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleVar)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar).Return(natgatewayRules, nil, nil)
		err := RunNatGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar).Return(natgatewayRules, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(&natgatewayRuleTestGet, nil, nil)
		err := RunNatGatewayRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, nil, nil)
		err := RunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, &testResponse, testNatGatewayRuleErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, nil, nil)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgRuleId), testNatGatewayRuleVar)
		cfg.Stdin = os.Stdin
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNatGatewayRulesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("rule", config.ArgCols), []string{"Name"})
	getNatGatewayRulesCols(core.GetGlobalFlagName("rule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNatGatewayRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("rule", config.ArgCols), []string{"Unknown"})
	getNatGatewayRulesCols(core.GetGlobalFlagName("rule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
