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
	testInputAlbForwardingRule = resources.ApplicationLoadBalancerForwardingRule{
		ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
			Properties: &ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
				Name:               &testAlbForwardingRuleVar,
				Protocol:           &testAlbForwardingRuleProtocol,
				ListenerIp:         &testAlbForwardingRuleVar,
				ListenerPort:       &testAlbForwardingRuleIntVar,
				ServerCertificates: &testAlbForwardingRuleServerCerts,
			}},
	}
	testAlbForwardingRulePropertiesNew = resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			Name:               &testAlbForwardingRuleNewVar,
			ListenerIp:         &testAlbForwardingRuleNewVar,
			ListenerPort:       &testAlbForwardingRuleNewIntVar,
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

func TestPreRunApplicationLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerForwardingRuleIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcApplicationLoadBalancerForwardingRuleIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(testAlbForwardingRules, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().ListForwardingRules(testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(testAlbForwardingRules, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleList(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(&testAlbForwardingRuleGet, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().GetForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(&testAlbForwardingRuleGet, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule).Return(&testAlbForwardingRuleGet, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule).Return(&testAlbForwardingRuleGet, &testResponse, nil)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule).Return(&testAlbForwardingRuleGet, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testAlbForwardingRuleProtocol)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().CreateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testInputAlbForwardingRule).Return(&testAlbForwardingRuleGet, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew).Return(&testAlbForwardingRuleUpdated, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew).Return(&testAlbForwardingRuleUpdated, nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerIp), testAlbForwardingRuleNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerPort), testAlbForwardingRuleNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerCertificates), testAlbForwardingRuleServerNewCerts)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().UpdateForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar, &testAlbForwardingRulePropertiesNew).Return(&testAlbForwardingRuleUpdated, nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(nil, testApplicationLoadBalancerForwardingRuleErr)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().DeleteForwardingRule(testAlbForwardingRuleVar, testAlbForwardingRuleVar, testAlbForwardingRuleVar).Return(nil, nil)
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerForwardingRuleDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testAlbForwardingRuleVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRuleId), testAlbForwardingRuleVar)
		cfg.Stdin = os.Stdin
		err := RunApplicationLoadBalancerForwardingRuleDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetAlbForwardingRulesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("forwardingrule", config.ArgCols), []string{"Name"})
	getAlbForwardingRulesCols(core.GetGlobalFlagName("forwardingrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetAlbForwardingRulesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("forwardingrule", config.ArgCols), []string{"Unknown"})
	getAlbForwardingRulesCols(core.GetGlobalFlagName("forwardingrule", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetApplicationLoadBalancerForwardingRulesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getAlbForwardingRulesIds(w, testAlbForwardingRuleVar, testAlbForwardingRuleVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
