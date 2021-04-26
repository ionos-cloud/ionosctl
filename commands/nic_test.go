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

func TestPreRunGlobalDcServerIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		err := PreRunGlobalDcServerIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerIdsValidateRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), "")
		err := PreRunGlobalDcServerIdsValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcServerIdsNicIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		err := PreRunGlobalDcServerIdsNicIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcServerIdsNicIdValidateRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), "")
		err := PreRunGlobalDcServerIdsNicIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		rm.Nic.EXPECT().List(testNicVar, testNicVar).Return(ns, nil, nil)
		err := RunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		rm.Nic.EXPECT().List(testNicVar, testNicVar).Return(ns, nil, testNicErr)
		err := RunNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		rm.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		rm.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicIps), ipsNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicIps), ipsNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicIps), ipsNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNicId)
		rm.Nic.EXPECT().Create(testNicVar, testNicVar, testNicVar, ipsNic, dhcpNic, lanNicId).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNewNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, nil)
		err := RunNicUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNewNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, testNicErr)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicName), testNicNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicDhcp), dhcpNewNic)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), lanNewNicId)
		rm.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties).Return(&nicNew, nil, nil)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, testNicErr)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testNicVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgNicId), testNicVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNicsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgCols), []string{"Name"})
	getNicsCols(builder.GetGlobalFlagName("nic", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNicsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgCols), []string{"Unknown"})
	getNicsCols(builder.GetGlobalFlagName("nic", config.ArgCols), w)
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
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgDataCenterId), testNicVar)
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgServerId), testNicVar)
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
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgDataCenterId), "")
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgServerId), "")
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
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgDataCenterId), testNicVar)
	viper.Set(builder.GetGlobalFlagName("nic", config.ArgLoadBalancerId), testNicVar)
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
