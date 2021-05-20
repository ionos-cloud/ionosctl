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
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	lanNicId    = int32(1)
	lanNewNicId = int32(2)
	dhcpNic     = false
	dhcpNewNic  = true
	ipsNic      = []string{"x.x.x.x"}
	n           = ionoscloud.Nic{
		Id: &testNicVar,
		Properties: &ionoscloud.NicProperties{
			Name:           &testNicVar,
			Lan:            &lanNicId,
			Dhcp:           &dhcpNic,
			Ips:            &ipsNic,
			FirewallActive: &dhcpNic,
			Mac:            &testNicVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
	}
	nicProperties = resources.NicProperties{
		NicProperties: ionoscloud.NicProperties{
			Name: &testNicNewVar,
			Dhcp: &dhcpNewNic,
			Lan:  &lanNewNicId,
		},
	}
	nicNew = resources.Nic{
		Nic: ionoscloud.Nic{
			Id: &testNicVar,
			Properties: &ionoscloud.NicProperties{
				Name:           nicProperties.NicProperties.Name,
				Lan:            nicProperties.NicProperties.Lan,
				Dhcp:           nicProperties.NicProperties.Dhcp,
				Ips:            &ipsNic,
				FirewallActive: &dhcpNic,
			},
		},
	}
	ns = resources.Nics{
		Nics: ionoscloud.Nics{
			Id:    &testNicVar,
			Items: &[]ionoscloud.Nic{n},
		},
	}
	balancedns = resources.BalancedNics{
		BalancedNics: ionoscloud.BalancedNics{
			Id:    &testNicVar,
			Items: &[]ionoscloud.Nic{n},
		},
	}
	testNicVar    = "test-nic"
	testNicNewVar = "test-new-nic"
	testNicErr    = errors.New("nic test: error occurred")
)

func TestPreRunGlobalDcServerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		err := PreRunGlobalDcServerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), "")
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), "")
		err := PreRunGlobalDcServerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerIdsNicId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		err := PreRunGlobalDcServerIdsNicId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerIdsNicIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), "")
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), "")
		err := PreRunGlobalDcServerIdsNicId(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		rm.Nic.EXPECT().List(testNicVar, testNicVar).Return(ns, nil, nil)
		err := RunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		rm.Nic.EXPECT().List(testNicVar, testNicVar).Return(ns, nil, testNicErr)
		err := RunNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		rm.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		rm.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, nil)
		err := RunNicUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, testNicErr)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, nil)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, testNicErr)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testNicVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNicsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgCols), []string{"Name"})
	getNicsCols(core.GetGlobalFlagName("nic", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNicsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgCols), []string{"Unknown"})
	getNicsCols(core.GetGlobalFlagName("nic", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetNicsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(core.GetGlobalFlagName("nic", config.ArgDataCenterId), testNicVar)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgServerId), testNicVar)
	getNicsIds(w, testNicVar, testNicVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetNicsIdsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(core.GetGlobalFlagName("nic", config.ArgDataCenterId), "")
	viper.Set(core.GetGlobalFlagName("nic", config.ArgServerId), "")
	getNicsIds(w, "", "")
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`404 Not Found`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetAttachedNicsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(core.GetGlobalFlagName("nic", config.ArgDataCenterId), testNicVar)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgLoadBalancerId), testNicVar)
	getAttachedNicsIds(w, testNicVar, testNicVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetAttachedNicsIdsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getAttachedNicsIds(w, "", "")
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`404 Not Found`)
	assert.True(t, re.Match(b.Bytes()))
}

// LoadBalancer Nic

func TestPreRunDcNicLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		err := PreRunDcNicLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNicLoadBalancerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), "")
		err := PreRunDcNicLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerNicAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerNicAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar).Return(balancedns, nil, nil)
		err := RunLoadBalancerNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar).Return(balancedns, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerNicGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		rm.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(&testResponse, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachAskForConfirm(t *testing.T) {
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar).Return(nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicDetachAskForConfirmErr(t *testing.T) {
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}
