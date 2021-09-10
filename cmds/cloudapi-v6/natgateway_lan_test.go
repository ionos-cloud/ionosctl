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
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	natgatewayLanTest = v6.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Properties: &ionoscloud.NatGatewayProperties{
				Name:      &testNatGatewayLanVar,
				PublicIps: &[]string{testNatGatewayLanVar},
				Lans:      &[]ionoscloud.NatGatewayLanProperties{natgatewayLanProperties.NatGatewayLanProperties},
			},
		},
	}
	natgatewayLanTestUpdated = v6.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Id: &testNatGatewayLanVar,
			Properties: &ionoscloud.NatGatewayProperties{
				Name:      &testNatGatewayLanVar,
				PublicIps: &[]string{testNatGatewayLanVar},
				Lans:      &[]ionoscloud.NatGatewayLanProperties{natgatewayLanProperties.NatGatewayLanProperties},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgatewayLanTestProper = v6.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &[]ionoscloud.NatGatewayLanProperties{
				natgatewayLanProperties.NatGatewayLanProperties,
				natgatewayLanNewProperties.NatGatewayLanProperties,
			},
		},
	}
	// Send empty struct to overwrite the existing one
	natgatewayLanTestRemove = v6.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &[]ionoscloud.NatGatewayLanProperties{
				natgatewayLanProperties.NatGatewayLanProperties,
			},
		},
	}
	natgatewayLanProperties = v6.NatGatewayLanProperties{
		NatGatewayLanProperties: ionoscloud.NatGatewayLanProperties{
			Id:         &testNatGatewayLanIntVar,
			GatewayIps: &[]string{testNatGatewayLanVar},
		},
	}
	natgatewayLanNewProperties = v6.NatGatewayLanProperties{
		NatGatewayLanProperties: ionoscloud.NatGatewayLanProperties{
			Id:         &testNatGatewayLanNewIntVar,
			GatewayIps: &[]string{testNatGatewayLanNewVar},
		},
	}
	testNatGatewayLanIntVar    = int32(1)
	testNatGatewayLanNewIntVar = int32(2)
	testNatGatewayLanVar       = "test-natgateway-lan"
	testNatGatewayLanNewVar    = "test-new-natgateway-lan"
	testNatGatewayLanErr       = errors.New("natgateway-lan test error")
)

func TestPreRunDcNatGatewayLanIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanVar)
		err := PreRunDcNatGatewayLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNatGatewayLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgCols), defaultNatGatewayLanCols)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNatGatewayLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper).Return(&natgatewayLanTestUpdated, nil, nil)
		err := RunNatGatewayLanAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanAddResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper).Return(&natgatewayLanTestUpdated, &testResponse, nil)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper).Return(&natgatewayLanTestUpdated, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayLanNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestProper).Return(&natgatewayLanTestUpdated, nil, nil)
		err := RunNatGatewayLanAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove).Return(&natgatewayLanTest, nil, testNatGatewayLanErr)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayLanRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayLanVar, testNatGatewayLanVar).Return(&natgatewayLanTest, nil, nil)
		rm.NatGateway.EXPECT().Update(testNatGatewayLanVar, testNatGatewayLanVar, natgatewayLanTestRemove).Return(&natgatewayLanTest, nil, nil)
		err := RunNatGatewayLanRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayLanRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testNatGatewayLanNewIntVar)
		cfg.Stdin = os.Stdin
		err := RunNatGatewayLanRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetNatGatewayLansCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Name"})
	getNatGatewayLansCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNatGatewayLansColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Unknown"})
	getNatGatewayLansCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
