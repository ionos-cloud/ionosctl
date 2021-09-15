package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testAlbRuleHttpRuleProperties = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{
				{
					Name:            &testAlbRuleHttpRuleVar,
					Type:            &testAlbRuleHttpRuleVar,
					TargetGroup:     &testAlbRuleHttpRuleVar,
					DropQuery:       &testAlbRuleHttpRuleBoolVar,
					Location:        &testAlbRuleHttpRuleVar,
					StatusCode:      &testAlbRuleHttpRuleIntVar,
					ResponseMessage: &testAlbRuleHttpRuleVar,
					ContentType:     &testAlbRuleHttpRuleVar,
					Conditions: &[]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
						{
							Type:      &testAlbRuleHttpRuleVar,
							Condition: &testAlbRuleHttpRuleVar,
							Key:       &testAlbRuleHttpRuleVar,
							Value:     &testAlbRuleHttpRuleVar,
							Negate:    &testAlbRuleHttpRuleBoolVar,
						},
					},
				},
			},
		},
	}
	testAlbRuleHttpRuleGet = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbRuleHttpRuleVar,
			Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{},
		},
	}
	testAlbRuleHttpRuleGetUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbRuleHttpRuleVar,
			Properties: &testAlbRuleHttpRuleProperties.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbRuleHttpRuleIntVar  = int32(1)
	testAlbRuleHttpRuleBoolVar = false
	testAlbRuleHttpRuleVar     = "test-rule-httprule"
	testAlbRuleHttpRuleErr     = errors.New("applicationloadbalancer-rule-httprule test error")
)

func TestPreRunApplicationLoadBalancerRuleHttpRule(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleIntVar)
		err := PreRunApplicationLoadBalancerRuleHttpRule(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerRuleHttpRuleErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunApplicationLoadBalancerRuleHttpRule(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleListGetHttpRulesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleListGetPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&resources.ApplicationLoadBalancerForwardingRule{}, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDropQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionKey), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleProperties).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDropQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionKey), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleProperties).Return(&testAlbRuleHttpRuleGetUpdated, &testResponse, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDropQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionKey), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleProperties).Return(&testAlbRuleHttpRuleGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDropQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionKey), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDropQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionKey), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleProperties).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}).Return(&testAlbRuleHttpRuleGet, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), "x.x.x.x")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemovePortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), "unknown type")
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar).Return(&testAlbRuleHttpRuleGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleVar)
		cfg.Stdin = os.Stdin
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetAlbRuleHttpRulesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("httprule", config.ArgCols), []string{"Name"})
	getAlbRuleHttpRulesCols(core.GetGlobalFlagName("httprule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetAlbRuleColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("httprule", config.ArgCols), []string{"Unknown"})
	getAlbRuleHttpRulesCols(core.GetGlobalFlagName("httprule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
