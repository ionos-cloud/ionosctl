package commands

import (
	"github.com/golang/mock/gomock"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testRule = resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Id: &testFirewallRuleVar,
			Properties: &ionoscloud.FirewallruleProperties{
				Name:           &testFirewallRuleVar,
				Protocol:       &testFirewallRuleProtocol,
				SourceMac:      &testFirewallRuleVar,
				SourceIp:       &testFirewallRuleVar,
				TargetIp:       &testFirewallRuleVar,
				IcmpCode:       &testFirewallRuleIntVar,
				IcmpType:       &testFirewallRuleIntVar,
				PortRangeStart: &testFirewallRulePortRangeStart,
				PortRangeEnd:   &testFirewallRulePortRangeEnd,
				Type:           &testFirewallRuleType,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testFirewallRuleState,
			},
		},
	}
	testFirewallRulesList = resources.FirewallRules{
		FirewallRules: ionoscloud.FirewallRules{
			Id: &testFirewallRuleVar,
			Items: &[]ionoscloud.FirewallRule{
				testRule.FirewallRule,
				testRule.FirewallRule,
			},
		},
	}
	testInputFirewallRuleProperties = resources.FirewallRuleProperties{
		FirewallruleProperties: ionoscloud.FirewallruleProperties{
			Name:           &testFirewallRuleVar,
			Protocol:       &testFirewallRuleProtocol,
			PortRangeStart: &testFirewallRulePortRangeStart,
			PortRangeEnd:   &testFirewallRulePortRangeEnd,
			SourceMac:      &testFirewallRuleVar,
			IcmpCode:       &testFirewallRuleIntVar,
			IcmpType:       &testFirewallRuleIntVar,
			SourceIp:       &testFirewallRuleVar,
			TargetIp:       &testFirewallRuleVar,
			Type:           &testFirewallRuleType,
		},
	}
	testInputFirewallRule = resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &testInputFirewallRuleProperties.FirewallruleProperties,
		},
	}
	testFirewallRules = resources.FirewallRules{
		FirewallRules: ionoscloud.FirewallRules{
			Id:    &testFirewallRuleVar,
			Items: &[]ionoscloud.FirewallRule{testRule.FirewallRule},
		},
	}
	testRequestIdVar = "f2354da4-83e3-4e92-9d23-f3cb1ffecc31"
	testResponse     = resources.Response{
		APIResponse: ionoscloud.APIResponse{
			Response: &http.Response{
				Header: map[string][]string{
					"Location": {"https://api.ionos.com/cloudapi/v6/requests/f2354da4-83e3-4e92-9d23-f3cb1ffecc31/status"},
				},
			},
			RequestTime: time.Duration(50),
		},
	}
	testResponseErr = resources.Response{
		APIResponse: ionoscloud.APIResponse{
			Response: &http.Response{
				Header: map[string][]string{
					"Location": {""},
				},
			},
			RequestTime: time.Duration(50),
		},
	}
	testFirewallRuleProtocol       = "TCP"
	testFirewallRuleType           = "INGRESS"
	testFirewallRuleState          = "AVAILABLE"
	testFirewallRulePortRangeStart = int32(2)
	testFirewallRulePortRangeEnd   = int32(2)
	testFirewallRuleIntVar         = int32(1)
	testFirewallRuleVar            = "test-firewall-rule"
	testFirewallRuleErr            = errors.New("firewall rule test error")
)

func TestFirewallruleCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(FirewallruleCmd())
	if ok := FirewallruleCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDcServerNicIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		err := PreRunDcServerNicIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerNicIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerNicIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunFirewallRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		err := PreRunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFirewallRuleListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunFirewallRuleListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunFirewallRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerNicIdsFRuleProtocol(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleVar)
		err := PreRunDcServerNicIdsFRuleProtocol(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerNicIdsFRuleProtocolErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerNicIdsFRuleProtocol(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerNicFRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		err := PreRunDcServerNicFRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerNicFRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerNicFRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFirewallRules, &testResponse, nil)
		err := RunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.FirewallRules{}, &testResponse, nil)
		err := RunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFirewallRules, nil, testFirewallRuleErr)
		err := RunFirewallRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Get(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testRule, &testResponse, nil)
		err := RunFirewallRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Get(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&testInputFirewallRule, &testResponse, nil)
		err := RunFirewallRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testInputFirewallRule, &testResponse, testFirewallRuleErr)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testInputFirewallRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testInputFirewallRule, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&testRule, &testResponse, nil)
		err := RunFirewallRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDestinationIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDirection), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testRule, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunFirewallRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&testResponse, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFirewallRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFirewallRulesList, nil, testFirewallRuleErr)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.FirewallRules{}, &testResponse, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.FirewallRules{FirewallRules: ionoscloud.FirewallRules{Items: &[]ionoscloud.FirewallRule{}}}, nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testFirewallRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testFirewallRuleErr)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testFirewallRuleErr)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetFirewallRulesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Name"})
	getFirewallRulesCols(core.GetGlobalFlagName("firewallrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetFirewallRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Unknown"})
	getFirewallRulesCols(core.GetGlobalFlagName("firewallrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
