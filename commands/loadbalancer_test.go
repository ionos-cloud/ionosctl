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
	testLoadbalancerVar    = "test-loadBalancer"
	testLoadbalancerNewVar = "test-new-loadBalancer"
	testLoadbalancerErr    = errors.New("loadBalancer test: error occurred")
)

func TestPreRunGlobalDcIdLoadbalancerIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		err := PreRunGlobalDcIdLoadbalancerIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdLoadbalancerIdValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		err := PreRunGlobalDcIdLoadbalancerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), "")
		err = PreRunGlobalDcIdLoadbalancerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgLoadbalancerId).Error())
	})
}

func TestRunLoadbalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, nil)
		err := RunLoadbalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().List(testLoadbalancerVar).Return(loadbs, nil, testLoadbalancerErr)
		err := RunLoadbalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Loadbalancer{loadb}, nil, nil)
		err := RunLoadbalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Loadbalancer{loadb}, nil, testLoadbalancerErr)
		err := RunLoadbalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{loadb}, nil, nil)
		err := RunLoadbalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerCreate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{loadb}, nil, testLoadbalancerErr)
		err := RunLoadbalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerCreate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancer)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer).Return(&resources.Loadbalancer{loadb}, nil, nil)
		err := RunLoadbalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadbalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerUpdate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, testLoadbalancerErr)
		err := RunLoadbalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerUpdate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerName), testLoadbalancerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerDhcp), dhcpLoadbalancerNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties).Return(&loadbalancerNew, nil, nil)
		err := RunLoadbalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadbalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerDelete_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, testLoadbalancerErr)
		err := RunLoadbalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerDelete_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadbalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadbalancerDelete_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadbalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadbalancerDelete_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLoadbalancerId), testLoadbalancerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunLoadbalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestLoadbalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("loadBalancer", config.ArgCols), []string{"Name"})
	getLoadbalancersCols(builder.GetGlobalFlagName("loadBalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLoadbalancersCols_Err(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("loadBalancer", config.ArgCols), []string{"Unknown"})
	getLoadbalancersCols(builder.GetGlobalFlagName("loadBalancer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetLoadbalancersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getLoadbalancersIds(w, "loadBalancer")
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
