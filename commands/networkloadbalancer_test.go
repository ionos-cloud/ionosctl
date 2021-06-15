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
	networkloadbalancerTest = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Properties: &ionoscloud.NetworkLoadBalancerProperties{
				Name:         &testNetworkLoadBalancerVar,
				Ips:          &[]string{testNetworkLoadBalancerVar},
				TargetLan:    &testNetworkLoadBalancerIntVar,
				ListenerLan:  &testNetworkLoadBalancerIntVar,
				LbPrivateIps: &[]string{testNetworkLoadBalancerVar},
			},
		},
	}
	networkloadbalancerTestGet = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Id:         &testNetworkLoadBalancerVar,
			Properties: networkloadbalancerTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	networkloadbalancers = resources.NetworkLoadBalancers{
		NetworkLoadBalancers: ionoscloud.NetworkLoadBalancers{
			Id:    &testNetworkLoadBalancerVar,
			Items: &[]ionoscloud.NetworkLoadBalancer{networkloadbalancerTest.NetworkLoadBalancer},
		},
	}
	networkloadbalancerProperties = resources.NetworkLoadBalancerProperties{
		NetworkLoadBalancerProperties: ionoscloud.NetworkLoadBalancerProperties{
			Name:         &testNetworkLoadBalancerNewVar,
			Ips:          &[]string{testNetworkLoadBalancerNewVar},
			TargetLan:    &testNetworkLoadBalancerNewIntVar,
			ListenerLan:  &testNetworkLoadBalancerNewIntVar,
			LbPrivateIps: &[]string{testNetworkLoadBalancerNewVar},
		},
	}
	networkloadbalancerNew = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Properties: &networkloadbalancerProperties.NetworkLoadBalancerProperties,
		},
	}
	testNetworkLoadBalancerIntVar    = int32(1)
	testNetworkLoadBalancerNewIntVar = int32(2)
	testNetworkLoadBalancerVar       = "test-networkloadbalancer"
	testNetworkLoadBalancerNewVar    = "test-new-networkloadbalancer"
	testNetworkLoadBalancerErr       = errors.New("networkloadbalancer test error")
)

func TestPreRunDcNetworkLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		err := PreRunDcNetworkLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar).Return(networkloadbalancers, nil, nil)
		err := RunNetworkLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar).Return(networkloadbalancers, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(&networkloadbalancerTestGet, nil, nil)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(&networkloadbalancerTestGet, nil, nil)
		rm.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(&networkloadbalancerTestGet, nil, nil)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(&networkloadbalancerTestGet, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(&networkloadbalancerTestGet, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest).Return(&networkloadbalancerTest, nil, nil)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest).Return(&networkloadbalancerTest, &testResponse, nil)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest).Return(&networkloadbalancerTest, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest).Return(&networkloadbalancerTest, nil, nil)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		rm.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties).Return(&networkloadbalancerNew, nil, nil)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		rm.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties).Return(&networkloadbalancerNew, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties).Return(&networkloadbalancerNew, nil, nil)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(nil, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(nil, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar).Return(nil, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		cfg.Stdin = os.Stdin
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNetworkLoadBalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), []string{"Name"})
	getNetworkLoadBalancersCols(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNetworkLoadBalancersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), []string{"Unknown"})
	getNetworkLoadBalancersCols(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetNetworkLoadBalancersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getNetworkLoadBalancersIds(w, testNetworkLoadBalancerVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
