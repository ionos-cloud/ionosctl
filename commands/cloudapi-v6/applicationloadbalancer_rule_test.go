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
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testInputAlbForwardingRule = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: compute.ApplicationLoadBalancerForwardingRule{
			Properties: compute.ApplicationLoadBalancerForwardingRuleProperties{
				Name:               &testAlbForwardingRuleVar,
				Protocol:           &testAlbForwardingRuleProtocol,
				ListenerIp:         &testAlbForwardingRuleVar,
				ListenerPort:       &testAlbForwardingRuleIntVar,
				ClientTimeout:      &testAlbForwardingRuleIntVar,
				ServerCertificates: &testAlbForwardingRuleServerCerts,
			}},
	}
	testAlbForwardingRulePropertiesNew = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: compute.ApplicationLoadBalancerForwardingRuleProperties{
			Name:               &testAlbForwardingRuleNewVar,
			ListenerIp:         &testAlbForwardingRuleNewVar,
			ListenerPort:       &testAlbForwardingRuleNewIntVar,
			ClientTimeout:      &testAlbForwardingRuleNewIntVar,
			ServerCertificates: &testAlbForwardingRuleServerNewCerts,
		},
	}
	testAlbForwardingRuleGet = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: compute.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbForwardingRuleVar,
			Properties: testInputAlbForwardingRule.Properties,
			Metadata: &compute.DatacenterElementMetadata{
				State: &testStateVar,
			},
		},
	}
	testAlbForwardingRuleUpdated = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: compute.ApplicationLoadBalancerForwardingRule{
			Id:         &testAlbForwardingRuleVar,
			Properties: testAlbForwardingRulePropertiesNew.ApplicationLoadBalancerForwardingRuleProperties,
		},
	}
	testAlbForwardingRules = resources.ApplicationLoadBalancerForwardingRules{
		ApplicationLoadBalancerForwardingRules: compute.ApplicationLoadBalancerForwardingRules{
			Items: &[]compute.ApplicationLoadBalancerForwardingRule{testAlbForwardingRuleGet.ApplicationLoadBalancerForwardingRule},
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		err := PreRunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		err := PreRunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerForwardingRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		err := PreRunDcApplicationLoadBalancerForwardingRuleIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerForwardingRuleIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcApplicationLoadBalancerForwardingRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgClientTimeout), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}
