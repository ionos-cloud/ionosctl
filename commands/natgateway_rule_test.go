package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	natgatewayRuleTest = v6.NatGatewayRule{
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
	natgatewayRuleTestGet = v6.NatGatewayRule{
		NatGatewayRule: ionoscloud.NatGatewayRule{
			Id:         &testNatGatewayRuleVar,
			Properties: natgatewayRuleTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgatewayRules = v6.NatGatewayRules{
		NatGatewayRules: ionoscloud.NatGatewayRules{
			Id:    &testNatGatewayRuleVar,
			Items: &[]ionoscloud.NatGatewayRule{natgatewayRuleTestGet.NatGatewayRule},
		},
	}
	natgatewayRuleProperties = v6.NatGatewayRuleProperties{
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
	natgatewayRuleNew = v6.NatGatewayRule{
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleVar)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar).Return(natgatewayRules, nil, nil)
		err := RunNatGatewayRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().ListRules(testNatGatewayRuleVar, testNatGatewayRuleVar).Return(natgatewayRules, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(&natgatewayRuleTestGet, nil, nil)
		err := RunNatGatewayRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().GetRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, nil, nil)
		err := RunNatGatewayRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, &testResponse, nil)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		rm.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().CreateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleTest).Return(&natgatewayRuleTestGet, nil, nil)
		err := RunNatGatewayRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, nil, nil)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		rm.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgProtocol), string(testNatGatewayRuleNewProtocol))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSourceSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetSubnet), testNatGatewayRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeStart), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPortRangeEnd), testNatGatewayRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().UpdateRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar, natgatewayRuleProperties).Return(&natgatewayRuleNew, nil, nil)
		err := RunNatGatewayRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		rm.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, testNatGatewayRuleErr)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.NatGateway.EXPECT().DeleteRule(testNatGatewayRuleVar, testNatGatewayRuleVar, testNatGatewayRuleVar).Return(nil, nil)
		err := RunNatGatewayRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNatGatewayRuleVar)
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

func TestGetNatGatewayRulesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getNatGatewayRulesIds(w, testNatGatewayRuleVar, testNatGatewayRuleVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
