package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testRuleTargetProperties = resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{
				{
					Ip:     &testNlbRuleTargetVar,
					Port:   &testNlbRuleTargetIntVar,
					Weight: &testNlbRuleTargetIntVar,
					HealthCheck: &ionoscloud.NetworkLoadBalancerForwardingRuleTargetHealthCheck{
						Check:         &testNlbRuleTargetBoolVar,
						CheckInterval: &testNlbRuleTargetIntVar,
						Maintenance:   &testNlbRuleTargetBoolVar,
					},
				},
			},
		},
	}
	testNlbRuleTargetGet = resources.NetworkLoadBalancerForwardingRule{
		NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
			Id:         &testNlbRuleTargetVar,
			Properties: &ionoscloud.NetworkLoadBalancerForwardingRuleProperties{},
		},
	}
	testNlbRuleTargetGetUpdated = resources.NetworkLoadBalancerForwardingRule{
		NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
			Id:         &testNlbRuleTargetVar,
			Properties: &testRuleTargetProperties.NetworkLoadBalancerForwardingRuleProperties,
		},
	}
	testTarget = ionoscloud.NetworkLoadBalancerForwardingRuleTarget{
		Ip:   &testNlbRuleTargetVar,
		Port: &testNlbRuleTargetIntVar,
	}
	testNlbRuleTarget = ionoscloud.NetworkLoadBalancerForwardingRule{
		Id: &testNlbRuleTargetVar,
		Properties: &ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Name:    &testNlbRuleTargetVar,
			Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{testTarget, testTarget},
		},
	}
	testNlbRuleTargetList = resources.NetworkLoadBalancerForwardingRules{
		NetworkLoadBalancerForwardingRules: ionoscloud.NetworkLoadBalancerForwardingRules{
			Id: &testNlbRuleTargetVar,
			Items: &[]ionoscloud.NetworkLoadBalancerForwardingRule{
				testNlbRuleTarget,
				testNlbRuleTarget,
			},
		},
	}
	testNlbRuleTargetIntVar  = int32(1)
	testNlbRuleTargetBoolVar = false
	testNlbRuleTargetVar     = "test-rule-target"
	testNlbRuleTargetErr     = errors.New("networkloadbalancer-rule-target test error")
)

func TestPreRunNetworkLoadBalancerRuleTarget(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			err := PreRunNetworkLoadBalancerRuleTarget(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunNetworkLoadBalancerRuleTargetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			err := PreRunNetworkLoadBalancerRuleTarget(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).
				Return(&testNlbRuleTargetGetUpdated, &testResponse, nil)
			err := RunNlbRuleTargetList(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNlbRuleTargetListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, testNlbRuleTargetErr)
			err := RunNlbRuleTargetList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetListGetTargetsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, nil)
			err := RunNlbRuleTargetList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetListGetPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&resources.NetworkLoadBalancerForwardingRule{}, nil, nil)
			err := RunNlbRuleTargetList(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testNlbRuleTargetBoolVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testNlbRuleTargetBoolVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar, &testRuleTargetProperties,
				testQueryParamOther,
			).
				Return(&testNlbRuleTargetGetUpdated, &testResponse, nil)
			err := RunNlbRuleTargetAdd(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNlbRuleTargetAddResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testNlbRuleTargetBoolVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testNlbRuleTargetBoolVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar, &testRuleTargetProperties,
				testQueryParamOther,
			).Return(&testNlbRuleTargetGetUpdated, &testResponse, testNlbRuleTargetErr)
			err := RunNlbRuleTargetAdd(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testNlbRuleTargetBoolVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testNlbRuleTargetBoolVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar, &testRuleTargetProperties,
				testQueryParamOther,
			).Return(&testNlbRuleTargetGetUpdated, nil, testNlbRuleTargetErr)
			err := RunNlbRuleTargetAdd(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testNlbRuleTargetBoolVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testNlbRuleTargetBoolVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, testNlbRuleTargetErr)
			err := RunNlbRuleTargetAdd(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testNlbRuleTargetBoolVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testNlbRuleTargetBoolVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGet, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar, &testRuleTargetProperties,
				testQueryParamOther,
			).Return(&testNlbRuleTargetGetUpdated, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunNlbRuleTargetAdd(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				&resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{},
					},
				},
				testQueryParamOther,
			).Return(&testNlbRuleTargetGet, &testResponse, nil)
			err := RunNlbRuleTargetRemove(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(constants.ArgVerbose, false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).
				Return(&testNlbRuleTargetGetUpdated, &testResponse, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				&resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{},
					},
				},
				testQueryParamOther,
			).Return(&testNlbRuleTargetGet, &testResponse, nil)
			err := RunNlbRuleTargetRemove(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				&resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{},
					},
				},
				testQueryParamOther,
			).Return(&testNlbRuleTargetGet, nil, testNlbRuleTargetErr)
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, testNlbRuleTargetErr)
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), "x.x.x.x")
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemovePortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), int32(2))
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), true)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				&resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{},
					},
				},
				testQueryParamOther,
			).Return(&testNlbRuleTargetGet, &testResponse, nil)
			rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(
				&testRequestStatus, nil, testRequestErr,
			)
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(constants.ArgQuiet, false)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().GetForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				gomock.AssignableToTypeOf(testQueryParamOther),
			).Return(&testNlbRuleTargetGetUpdated, nil, nil)
			rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().UpdateForwardingRule(
				testNlbRuleTargetVar, testNlbRuleTargetVar, testNlbRuleTargetVar,
				&resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget{},
					},
				},
				testQueryParamOther,
			).Return(&testNlbRuleTargetGet, nil, nil)
			err := RunNlbRuleTargetRemove(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestRunNlbRuleTargetRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(constants.ArgQuiet, false)
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgForce), false)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testNlbRuleTargetVar)
			viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testNlbRuleTargetIntVar)
			cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
			err := RunNlbRuleTargetRemove(cfg)
			assert.Error(t, err)
		},
	)
}
