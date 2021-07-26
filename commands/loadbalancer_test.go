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
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	dhcpLoadbalancer    = true
	dhcpLoadbalancerNew = false
	loadb               = ionoscloud.Loadbalancer{
		Id: &testLoadbalancerVar,
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &testLoadbalancerVar,
			Dhcp: &dhcpLoadbalancer,
			Ip:   &testLoadbalancerVar,
		},
	}
	loadbs = v5.Loadbalancers{
		Loadbalancers: ionoscloud.Loadbalancers{
			Id:    &testLoadbalancerVar,
			Items: &[]ionoscloud.Loadbalancer{loadb},
		},
	}
	loadbalancerProperties = v5.LoadbalancerProperties{
		LoadbalancerProperties: ionoscloud.LoadbalancerProperties{
			Name: &testLoadbalancerNewVar,
			Dhcp: &dhcpLoadbalancerNew,
			Ip:   &testLoadbalancerNewVar,
		},
	}
	loadbalancerNew = v5.Loadbalancer{
		Loadbalancer: ionoscloud.Loadbalancer{
			Id:         &testLoadbalancerVar,
			Properties: &loadbalancerProperties.LoadbalancerProperties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	testLoadbalancerVar    = "test-loadbalancer"
	testLoadbalancerNewVar = "test-new-loadbalancer"
	testLoadbalancerErr    = errors.New("loadbalancer test: error occurred")
)

func TestPreRunDcLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		err := PreRunDcLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLoadBalancerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, nil)
		err := RunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, testLoadbalancerErr)
		err := RunLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&v5.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&v5.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&v5.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&v5.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&v5.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, testLoadbalancerErr)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, &testResponse, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, testLoadbalancerErr)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestLoadbalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("loadbalancer", config.ArgCols), []string{"Name"})
	getLoadbalancersCols(core.GetGlobalFlagName("loadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLoadbalancersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("loadbalancer", config.ArgCols), []string{"Unknown"})
	getLoadbalancersCols(core.GetGlobalFlagName("loadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetLoadbalancersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getLoadbalancersIds(w, "loadbalancer")
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
