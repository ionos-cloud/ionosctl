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
	testInputAlbForwardingRule = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
				Name:               &testAlbForwardingRuleVar,
				Protocol:           &testAlbForwardingRuleProtocol,
				ListenerIp:         &testAlbForwardingRuleVar,
				ListenerPort:       &testAlbForwardingRuleIntVar,
				ClientTimeout:      &testAlbForwardingRuleIntVar,
				ServerCertificates: &testAlbForwardingRuleServerCerts,
			}},
	}
	testAlbForwardingRulePropertiesNew = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			Name:               &testAlbForwardingRuleNewVar,
			ListenerIp:         &testAlbForwardingRuleNewVar,
			ListenerPort:       &testAlbForwardingRuleNewIntVar,
			ClientTimeout:      &testAlbForwardingRuleNewIntVar,
			ServerCertificates: &testAlbForwardingRuleServerNewCerts,
		},
	}
	testAlbForwardingRuleGet = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbForwardingRuleVar,
			Properties: testInputAlbForwardingRule.Properties,
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	testAlbForwardingRuleUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbForwardingRuleVar,
			Properties: &testAlbForwardingRulePropertiesNew.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbForwardingRules = resources.ApplicationLoadBalancerForwardingRules{
		ApplicationLoadBalancerForwardingRules: ionoscloud.ApplicationLoadBalancerForwardingRules{
			Items: &[]ionoscloud.ApplicationLoadBalancerForwardingRule{testAlbForwardingRuleGet.ApplicationLoadBalancerForwardingRule},
		},
	}
	testAlbForwardingRuleIntVar                  = int32(1)
	testAlbForwardingRuleNewIntVar               = int32(2)
	testAlbForwardingRuleProtocol                = "HTTP"
	testAlbForwardingRuleVar                     = "test-forwardingrule"
	testAlbForwardingRuleNewVar                  = "test-new-forwardingrule"
	testAlbForwardingRuleServerCerts             = []string{testAlbForwardingRuleVar}
	testAlbForwardingRuleServerNewCerts          = []string{testAlbForwardingRuleNewVar}
	testApplicationLoadBalancerForwardingRuleErr = errors.New("applicationloadbalancer-forwardingrule test error")
)

func TestApplicationLoadBalancerRuleCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ApplicationLoadBalancerRuleCmd())
	if ok := ApplicationLoadBalancerRuleCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunApplicationLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleIntVar)
		err := PreRunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerForwardingRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		err := PreRunDcApplicationLoadBalancerForwardingRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerForwardingRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcApplicationLoadBalancerForwardingRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, &testResponse, nil)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, &testResponse, nil)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, &testResponse, nil)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testAlbForwardingRuleGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew, testQueryParamOther).Return(&testAlbForwardingRuleUpdated, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew, testQueryParamOther).Return(&testAlbForwardingRuleUpdated, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew, testQueryParamOther).Return(&testAlbForwardingRuleUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testApplicationLoadBalancerForwardingRuleErr)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testAlbForwardingRules, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testApplicationLoadBalancerForwardingRuleErr)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRuleId), testAlbForwardingRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}
