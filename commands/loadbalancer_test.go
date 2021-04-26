package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
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
	loadbs = resources.Loadbalancers{
		Loadbalancers: ionoscloud.Loadbalancers{
			Id:    &testLoadbalancerVar,
			Items: &[]ionoscloud.Loadbalancer{loadb},
		},
	}
	loadbalancerProperties = resources.LoadbalancerProperties{
		LoadbalancerProperties: ionoscloud.LoadbalancerProperties{
			Name: &testLoadbalancerNewVar,
			Dhcp: &dhcpLoadbalancerNew,
		},
	}
	loadbalancerNew = resources.Loadbalancer{
		Loadbalancer: ionoscloud.Loadbalancer{
			Id: &testLoadbalancerVar,
			Properties: &ionoscloud.LoadbalancerProperties{
				Name: loadbalancerProperties.LoadbalancerProperties.Name,
				Dhcp: &dhcpLoadbalancerNew,
			},
		},
	}
	testLoadbalancerVar    = "test-loadbalancer"
	testLoadbalancerNewVar = "test-new-loadbalancer"
	testLoadbalancerErr    = errors.New("loadbalancer test: error occurred")
)

func TestPreRunGlobalDcIdLoadBalancerIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		err := PreRunGlobalDcIdLoadBalancerIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdLoadBalancerIdValidateRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), "")
		err := PreRunGlobalDcIdLoadBalancerIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, nil)
		err := RunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, testLoadbalancerErr)
		err := RunLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, nil)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, testLoadbalancerErr)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, testLoadbalancerErr)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestLoadbalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("loadbalancer", config.ArgCols), []string{"Name"})
	getLoadbalancersCols(builder.GetGlobalFlagName("loadbalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLoadbalancersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("loadbalancer", config.ArgCols), []string{"Unknown"})
	getLoadbalancersCols(builder.GetGlobalFlagName("loadbalancer", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getLoadbalancersIds(w, "loadbalancer")
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestPreRunGlobalDcIdNicLoadBalancerIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		err := PreRunGlobalDcIdNicLoadBalancerIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdNicLoadBalancerIdsValidateRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), "")
		err := PreRunGlobalDcIdNicLoadBalancerIdsValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerAttachNic(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerAttachNic(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerAttachNicErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerAttachNic(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerAttachNicWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerAttachNic(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerListNics(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar).Return(balancedns, nil, nil)
		err := RunLoadBalancerListNics(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListNicsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar).Return(balancedns, nil, testLoadbalancerErr)
		err := RunLoadBalancerListNics(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerGetNic(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerGetNic(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerGetNicErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerGetNic(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDetachNic(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDetachNic(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDetachNicErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&testResponse, nil)
		err := RunLoadBalancerDetachNic(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDetachNicWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDetachNic(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDetachNicAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerDetachNic(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDetachNicAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunLoadBalancerDetachNic(cfg)
		assert.Error(t, err)
	})
}
