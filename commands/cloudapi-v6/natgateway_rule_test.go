package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		err := PreRunNATGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunNATGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunNATGatewayRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleVar)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcNatGatewayRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNatGatewayRuleVar)
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
	viper.Set(core.GetGlobalFlagName("rule", constants.ArgCols), []string{"Name"})
	getNatGatewayRulesCols(core.GetGlobalFlagName("rule", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNatGatewayRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("rule", constants.ArgCols), []string{"Unknown"})
	getNatGatewayRulesCols(core.GetGlobalFlagName("rule", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
