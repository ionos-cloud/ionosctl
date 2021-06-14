package commands

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
	testResponse = resources.Response{
		APIResponse: ionoscloud.APIResponse{
			Response: &http.Response{
				Header: map[string][]string{
					"location": {"https://api.ionos.com/cloudapi/v6/create/resource/status"},
				},
			},
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

func TestPreRunGlobalDcServerNicIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcServerNicIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleProtocol(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIdsFRuleProtocol(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleProtocolErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcServerNicIdsFRuleProtocol(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIdsFRuleId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcServerNicIdsFRuleId(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(testFirewallRules, nil, nil)
		err := RunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(testFirewallRules, nil, testFirewallRuleErr)
		err := RunFirewallRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().Get(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().Get(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(&testRule, nil, nil)
		err := RunFirewallRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, nil, nil)
		err := RunFirewallRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, &testResponse, nil)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, nil, nil)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties).Return(&testRule, nil, nil)
		err := RunFirewallRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), testFirewallRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testFirewallRulePortRangeEnd)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpCode), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIcmpType), testFirewallRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceMac), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetIp), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testFirewallRuleType)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties).Return(&testRule, nil, nil)
		err := RunFirewallRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, testFirewallRuleErr)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testFirewallRuleVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgNicId), testFirewallRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirewallRuleId), testFirewallRuleVar)
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

func TestGetFirewallRulesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getFirewallRulesIds(w, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
