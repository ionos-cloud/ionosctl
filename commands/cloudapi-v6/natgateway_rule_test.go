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
	natgatewayRulesList = resources.NatGatewayRules{
		NatGatewayRules: ionoscloud.NatGatewayRules{
			Id: &testNatGatewayRuleVar,
			Items: &[]ionoscloud.NatGatewayRule{
				natgatewayRuleTestGet.NatGatewayRule,
				natgatewayRuleTestGet.NatGatewayRule,
			},
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

func TestPreRunNatGatewayRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		err := PreRunNATGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunNATGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunNATGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleVar)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(natgatewayRules, &testResponse, nil)
		err := RunNatGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.NatGatewayRules{}, &testResponse, nil)
		err := RunNatGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(natgatewayRules, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&natgatewayRuleTestGet, &testResponse, nil)
		err := RunNatGatewayRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&natgatewayRuleTestGet, &testResponse, nil)
		err := RunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleTestGet, &testResponse, testNatGatewayRuleErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleTestGet, &testResponse, nil)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&natgatewayRuleNew, &testResponse, nil)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleNew, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPortRangeEnd), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&natgatewayRuleNew, &testResponse, nil)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(natgatewayRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}
func TestRunNatGatewayRuleDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(natgatewayRulesList, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}
func TestRunNatGatewayRuleDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.NatGatewayRules{}, &testResponse, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}
func TestRunNatGatewayRuleDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.NatGatewayRules{NatGatewayRules: ionoscloud.NatGatewayRules{Items: &[]ionoscloud.NatGatewayRule{}}}, &testResponse, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(natgatewayRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testNatGatewayRuleErr)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNatGatewayRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}
