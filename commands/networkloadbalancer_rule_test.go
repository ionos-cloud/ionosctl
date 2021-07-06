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
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testInputForwardingRule = resources.NetworkLoadBalancerForwardingRule{
		NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
			Properties: &ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
				Name:         &testNlbForwardingRuleVar,
				Algorithm:    &testNlbForwardingRuleAlgorithm,
				Protocol:     &testNlbForwardingRuleProtocol,
				ListenerIp:   &testNlbForwardingRuleVar,
				ListenerPort: &testNlbForwardingRuleIntVar,
				HealthCheck: &ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{
					ClientTimeout:  &testNlbForwardingRuleIntVar,
					ConnectTimeout: &testNlbForwardingRuleIntVar,
					TargetTimeout:  &testNlbForwardingRuleIntVar,
					Retries:        &testNlbForwardingRuleIntVar,
				},
			}},
	}
	testNlbForwardingRulePropertiesNew = resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Name:         &testNlbForwardingRuleNewVar,
			ListenerIp:   &testNlbForwardingRuleNewVar,
			ListenerPort: &testNlbForwardingRuleNewIntVar,
			HealthCheck: &ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{
				ClientTimeout:  &testNlbForwardingRuleNewIntVar,
				ConnectTimeout: &testNlbForwardingRuleNewIntVar,
				TargetTimeout:  &testNlbForwardingRuleNewIntVar,
				Retries:        &testNlbForwardingRuleNewIntVar,
			},
		},
	}
	testNlbForwardingRuleGet = resources.NetworkLoadBalancerForwardingRule{
		NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
			Id:         &testNlbForwardingRuleVar,
			Properties: testInputForwardingRule.Properties,
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	testNlbForwardingRuleUpdated = resources.NetworkLoadBalancerForwardingRule{
		NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
			Id:         &testNlbForwardingRuleVar,
			Properties: &testNlbForwardingRulePropertiesNew.NetworkLoadBalancerForwardingRuleProperties,
		},
	}
	testNlbForwardingRules = resources.NetworkLoadBalancerForwardingRules{
		NetworkLoadBalancerForwardingRules: ionoscloud.NetworkLoadBalancerForwardingRules{
			Items: &[]ionoscloud.NetworkLoadBalancerForwardingRule{testNlbForwardingRuleGet.NetworkLoadBalancerForwardingRule},
		},
	}
	testNlbForwardingRuleIntVar              = int32(1)
	testNlbForwardingRuleNewIntVar           = int32(2)
	testNlbForwardingRuleProtocol            = "TCP"
	testNlbForwardingRuleAlgorithm           = "ROUND_ROBIN"
	testNlbForwardingRuleVar                 = "test-forwardingrule"
	testNlbForwardingRuleNewVar              = "test-new-forwardingrule"
	testNetworkLoadBalancerForwardingRuleErr = errors.New("networkloadbalancer-forwardingrule test error")
)

func TestPreRunNetworkLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleIntVar)
		err := PreRunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerForwardingRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		err := PreRunDcNetworkLoadBalancerForwardingRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerForwardingRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerForwardingRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(testNlbForwardingRules, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(testNlbForwardingRules, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().GetForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(&testNlbForwardingRuleGet, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().GetForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(&testNlbForwardingRuleGet, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleIntVar)
		rm.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testInputForwardingRule).Return(&testNlbForwardingRuleGet, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleIntVar)
		rm.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testInputForwardingRule).Return(&testNlbForwardingRuleGet, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleIntVar)
		rm.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testInputForwardingRule).Return(&testNlbForwardingRuleGet, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleIntVar)
		rm.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testInputForwardingRule).Return(&testNlbForwardingRuleGet, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleNewIntVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew).Return(&testNlbForwardingRuleUpdated, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleNewIntVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew).Return(&testNlbForwardingRuleUpdated, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRetries), testNlbForwardingRuleNewIntVar)
		rm.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew).Return(&testNlbForwardingRuleUpdated, nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		rm.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar).Return(nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRuleId), testNlbForwardingRuleVar)
		cfg.Stdin = os.Stdin
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetForwardingRulesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("forwardingrule", config.ArgCols), []string{"Name"})
	getForwardingRulesCols(core.GetGlobalFlagName("forwardingrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetForwardingRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("forwardingrule", config.ArgCols), []string{"Unknown"})
	getForwardingRulesCols(core.GetGlobalFlagName("forwardingrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetNetworkLoadBalancerForwardingRulesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getForwardingRulesIds(w, testNlbForwardingRuleVar, testNlbForwardingRuleVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
