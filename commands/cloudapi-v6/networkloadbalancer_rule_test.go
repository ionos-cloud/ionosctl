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
	testNlbForwardingRulesList = resources.NetworkLoadBalancerForwardingRules{
		NetworkLoadBalancerForwardingRules: ionoscloud.NetworkLoadBalancerForwardingRules{
			Items: &[]ionoscloud.NetworkLoadBalancerForwardingRule{
				testNlbForwardingRuleGet.NetworkLoadBalancerForwardingRule,
				testNlbForwardingRuleGet.NetworkLoadBalancerForwardingRule,
			},
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

func TestPreRunNetworkLoadBalancerRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		err := PreRunNetworkLoadBalancerRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerRuleListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunNetworkLoadBalancerRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunNetworkLoadBalancerRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleIntVar)
		err := PreRunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerForwardingRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		err := PreRunDcNetworkLoadBalancerForwardingRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerForwardingRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerForwardingRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(testNlbForwardingRules, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(resources.NetworkLoadBalancerForwardingRules{}, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testNlbForwardingRules, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&testNlbForwardingRuleGet, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testNlbForwardingRuleGet, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testInputForwardingRule, testQueryParamOther,
		).
			Return(&testNlbForwardingRuleGet, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testInputForwardingRule, testQueryParamOther,
		).Return(&testNlbForwardingRuleGet, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testInputForwardingRule, testQueryParamOther,
		).Return(&testNlbForwardingRuleGet, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testNlbForwardingRuleAlgorithm)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().CreateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testInputForwardingRule, testQueryParamOther,
		).Return(&testNlbForwardingRuleGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleNewIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew, testQueryParamOther,
		).
			Return(&testNlbForwardingRuleUpdated, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleNewIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew, testQueryParamOther,
		).Return(&testNlbForwardingRuleUpdated, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testNlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagConnectionTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetTimeout), testNlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testNlbForwardingRuleNewIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, &testNlbForwardingRulePropertiesNew, testQueryParamOther,
		).Return(&testNlbForwardingRuleUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar, testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).
			Return(&testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(testNlbForwardingRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).
			Return(&testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).
			Return(&testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(testNlbForwardingRulesList, nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(resources.NetworkLoadBalancerForwardingRules{}, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(resources.NetworkLoadBalancerForwardingRules{NetworkLoadBalancerForwardingRules: ionoscloud.NetworkLoadBalancerForwardingRules{Items: &[]ionoscloud.NetworkLoadBalancerForwardingRule{}}}, &testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().ListForwardingRules(testNlbForwardingRuleVar, testNlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).
			Return(testNlbForwardingRulesList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).
			Return(&testResponse, testNetworkLoadBalancerForwardingRuleErr)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).
			Return(&testResponse, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).Return(nil, testNetworkLoadBalancerForwardingRuleErr)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().DeleteForwardingRule(testNlbForwardingRuleVar,
			testNlbForwardingRuleVar, testNlbForwardingRuleVar, testQueryParamOther,
		).Return(nil, nil)
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerForwardingRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNetworkLoadBalancerId), testNlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testNlbForwardingRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunNetworkLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}
