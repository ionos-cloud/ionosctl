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
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	applicationloadbalancerTest = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Properties: &ionoscloud.ApplicationLoadBalancerProperties{
				Name:         &testApplicationLoadBalancerVar,
				Ips:          &[]string{testApplicationLoadBalancerVar},
				TargetLan:    &testApplicationLoadBalancerIntVar,
				ListenerLan:  &testApplicationLoadBalancerIntVar,
				LbPrivateIps: &[]string{testApplicationLoadBalancerVar},
			},
		},
	}
	applicationloadbalancerTestGet = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Id:         &testApplicationLoadBalancerVar,
			Properties: applicationloadbalancerTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	applicationloadbalancers = resources.ApplicationLoadBalancers{
		ApplicationLoadBalancers: ionoscloud.ApplicationLoadBalancers{
			Id:    &testApplicationLoadBalancerVar,
			Items: &[]ionoscloud.ApplicationLoadBalancer{applicationloadbalancerTestGet.ApplicationLoadBalancer},
		},
	}
	applicationloadbalancerProperties = resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: ionoscloud.ApplicationLoadBalancerProperties{
			Name:         &testApplicationLoadBalancerNewVar,
			Ips:          &[]string{testApplicationLoadBalancerNewVar},
			TargetLan:    &testApplicationLoadBalancerNewIntVar,
			ListenerLan:  &testApplicationLoadBalancerNewIntVar,
			LbPrivateIps: &[]string{testApplicationLoadBalancerNewVar},
		},
	}
	applicationloadbalancerNew = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Properties: &applicationloadbalancerProperties.ApplicationLoadBalancerProperties,
		},
	}
	testApplicationLoadBalancerIntVar    = int32(1)
	testApplicationLoadBalancerNewIntVar = int32(2)
	testApplicationLoadBalancerVar       = "test-applicationloadbalancer"
	testApplicationLoadBalancerNewVar    = "test-new-applicationloadbalancer"
	testApplicationLoadBalancerErr       = errors.New("applicationloadbalancer test error")
)

func TestApplicationLoadBalancerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ApplicationLoadBalancerCmd())
	if ok := ApplicationLoadBalancerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDcApplicationLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		err := PreRunDcApplicationLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcApplicationLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunApplicationLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		err := PreRunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		err := PreRunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, &testResponse, nil)
		err := RunApplicationLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, nil, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, &testResponse, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, nil, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&applicationloadbalancerTestGet, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest).Return(&applicationloadbalancerTest, nil, nil)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest).Return(&applicationloadbalancerTest, &testResponse, nil)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest).Return(&applicationloadbalancerTest, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest).Return(&applicationloadbalancerTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties).Return(&applicationloadbalancerNew, nil, nil)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties).Return(&applicationloadbalancerNew, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties).Return(&applicationloadbalancerNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, &testResponse, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, &testResponse, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&testResponse, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, &testResponse, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		cfg.Stdin = os.Stdin
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar).Return(applicationloadbalancers, &testResponse, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		cfg.Stdin = os.Stdin
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetApplicationLoadBalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("applicationloadbalancer", config.ArgCols), []string{"Name"})
	getApplicationLoadBalancersCols(core.GetGlobalFlagName("applicationloadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetApplicationLoadBalancersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("applicationloadbalancer", config.ArgCols), []string{"Unknown"})
	getApplicationLoadBalancersCols(core.GetGlobalFlagName("applicationloadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
