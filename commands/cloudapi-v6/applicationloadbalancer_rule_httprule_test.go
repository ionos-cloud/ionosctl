package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testAlbRuleHttpRuleForwardProperties = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{
				{
					Name:        &testAlbRuleHttpRuleVar,
					Type:        &testAlbRuleHttpRuleForwardTypeVar,
					TargetGroup: &testAlbRuleHttpRuleVar,
					Conditions: &[]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
						{
							Type:      &testAlbRuleHttpRuleVar,
							Condition: &testAlbRuleHttpRuleVar,
							Value:     &testAlbRuleHttpRuleVar,
							Negate:    &testAlbRuleHttpRuleBoolVar,
						},
					},
				},
			},
		},
	}
	testAlbRuleHttpRuleStaticProperties = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{
				{
					Name:            &testAlbRuleHttpRuleVar,
					Type:            &testAlbRuleHttpRuleStaticTypeVar,
					StatusCode:      &testAlbRuleHttpRuleIntVar,
					ResponseMessage: &testAlbRuleHttpRuleVar,
					ContentType:     &testAlbRuleHttpRuleVar,
					Conditions: &[]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
						{
							Type:      &testAlbRuleHttpRuleVar,
							Condition: &testAlbRuleHttpRuleVar,
							Value:     &testAlbRuleHttpRuleVar,
							Negate:    &testAlbRuleHttpRuleBoolVar,
						},
					},
				},
			},
		},
	}
	testAlbRuleHttpRuleRedirectProperties = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{
				{
					Name:       &testAlbRuleHttpRuleVar,
					Type:       &testAlbRuleHttpRuleRedirectTypeVar,
					StatusCode: &testAlbRuleHttpRuleIntVar,
					Location:   &testAlbRuleHttpRuleVar,
					DropQuery:  &testAlbRuleHttpRuleBoolVar,
					Conditions: &[]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
						{
							Type:      &testAlbRuleHttpRuleVar,
							Condition: &testAlbRuleHttpRuleVar,
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
	testAlbRuleHttpRuleForwardGetUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbRuleHttpRuleVar,
			Properties: &testAlbRuleHttpRuleForwardProperties.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbRuleHttpRuleStaticGetUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbRuleHttpRuleVar,
			Properties: &testAlbRuleHttpRuleStaticProperties.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbRuleHttpRuleRedirectGetUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbRuleHttpRuleVar,
			Properties: &testAlbRuleHttpRuleRedirectProperties.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbRuleHttpRuleIntVar          = int32(1)
	testAlbRuleHttpRuleBoolVar         = false
	testAlbRuleHttpRuleForwardTypeVar  = "FORWARD"
	testAlbRuleHttpRuleStaticTypeVar   = "STATIC"
	testAlbRuleHttpRuleRedirectTypeVar = "REDIRECT"
	testAlbRuleHttpRuleVar             = "test-rule-httprule"
	testAlbRuleHttpRuleErr             = errors.New("applicationloadbalancer-rule-httprule test error")
)

func TestAlbRuleHttpRuleCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(AlbRuleHttpRuleCmd())
	if ok := AlbRuleHttpRuleCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunApplicationLoadBalancerRuleHttpRule(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunApplicationLoadBalancerRuleHttpRule(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunApplicationLoadBalancerRuleHttpRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		err := PreRunApplicationLoadBalancerRuleHttpRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerRuleHttpRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		err := PreRunApplicationLoadBalancerRuleHttpRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerRuleHttpRuleDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		err := PreRunApplicationLoadBalancerRuleHttpRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, &testResponse, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleListGetHttpRulesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleListGetPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.ApplicationLoadBalancerForwardingRule{}, nil, nil)
		err := RunAlbRuleHttpRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddForward(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleForwardTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleForwardProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddStatic(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleStaticTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMessage), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleStaticProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleStaticGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddRedirect(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleRedirectTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgQuery), testAlbRuleHttpRuleBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgStatusCode), testAlbRuleHttpRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgContentType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleRedirectProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleRedirectGetUpdated, nil, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleForwardTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleForwardProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleForwardGetUpdated, &testResponse, nil)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleForwardTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleForwardProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleForwardTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testAlbRuleHttpRuleForwardTypeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCondition), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionType), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConditionValue), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testAlbRuleHttpRuleBoolVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, &testAlbRuleHttpRuleForwardProperties, testQueryParamOther).Return(&testAlbRuleHttpRuleForwardGetUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunAlbRuleHttpRuleAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}, testQueryParamOther).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}, testQueryParamOther).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}, testQueryParamOther).Return(&testAlbRuleHttpRuleGet, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, testAlbRuleHttpRuleErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}, testQueryParamOther).Return(&testAlbRuleHttpRuleGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbRuleHttpRuleForwardGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar, testAlbRuleHttpRuleVar,
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
					HttpRules: &[]ionoscloud.ApplicationLoadBalancerHttpRule{},
				},
			}, testQueryParamOther).Return(&testAlbRuleHttpRuleGet, nil, nil)
		err := RunAlbRuleHttpRuleRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunAlbRuleHttpRuleRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbRuleHttpRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbRuleHttpRuleVar)
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
	viper.Set(core.GetGlobalFlagName("httprule", constants.ArgCols), []string{"Name"})
	getAlbRuleHttpRulesCols(core.GetGlobalFlagName("httprule", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetAlbRuleColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("httprule", constants.ArgCols), []string{"Unknown"})
	getAlbRuleHttpRulesCols(core.GetGlobalFlagName("httprule", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
