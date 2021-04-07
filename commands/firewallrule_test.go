package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
				PortRangeStart: &testFirewallRulePortRangeStart,
				PortRangeEnd:   &testFirewallRulePortRangeEnd,
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
	testFirewallRuleProtocol       = "TCP"
	testFirewallRuleState          = "AVAILABLE"
	testFirewallRulePortRangeStart = int32(2)
	testFirewallRulePortRangeEnd   = int32(2)
	testFirewallRuleVar            = "test-firewall-rule"
	testFirewallRuleErr            = errors.New("firewall rule test error")
)

func TestPreRunGlobalDcServerNicIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), "")
		err := PreRunGlobalDcServerNicIdsValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleProtocolValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIdsFRuleProtocolValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleProtocolValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), "")
		err := PreRunGlobalDcServerNicIdsFRuleProtocolValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		err := PreRunGlobalDcServerNicIdsFRuleIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerNicIdsFRuleIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), "")
		err := PreRunGlobalDcServerNicIdsFRuleIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(testFirewallRules, nil, nil)
		err := RunFirewallRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().List(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(testFirewallRules, nil, testFirewallRuleErr)
		err := RunFirewallRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		rm.FirewallRule.EXPECT().Get(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleName), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), testFirewallRuleProtocol)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStop), testFirewallRulePortRangeEnd)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, nil, nil)
		err := RunFirewallRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleName), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), testFirewallRuleProtocol)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStop), testFirewallRulePortRangeEnd)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Create(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRule).Return(&testInputFirewallRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleName), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), testFirewallRuleProtocol)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStop), testFirewallRulePortRangeEnd)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties).Return(&testRule, nil, nil)
		err := RunFirewallRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleName), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleProtocol), testFirewallRuleProtocol)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStart), testFirewallRulePortRangeStart)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRulePortRangeStop), testFirewallRulePortRangeEnd)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Update(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testInputFirewallRuleProperties).Return(&testRule, nil, testFirewallRuleErr)
		err := RunFirewallRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, testFirewallRuleErr)
		err := RunFirewallRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.FirewallRule.EXPECT().Delete(testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar, testFirewallRuleVar).Return(nil, nil)
		err := RunFirewallRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunFirewallRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testFirewallRuleVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgNicId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgFirewallRuleId), testFirewallRuleVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
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
	viper.Set(builder.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Name"})
	getFirewallRulesCols(builder.GetGlobalFlagName("firewallrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetFirewallRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("firewallrule", config.ArgCols), []string{"Unknown"})
	getFirewallRulesCols(builder.GetGlobalFlagName("firewallrule", config.ArgCols), w)
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
