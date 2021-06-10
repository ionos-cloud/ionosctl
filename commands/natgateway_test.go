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
	natgatewayTest = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Properties: &ionoscloud.NatGatewayProperties{
				Name:      &testNatGatewayVar,
				PublicIps: &[]string{testNatGatewayVar},
			},
		},
	}
	natgatewayTestGet = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Id:         &testNatGatewayVar,
			Properties: natgatewayTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	natgateways = resources.NatGateways{
		NatGateways: ionoscloud.NatGateways{
			Id:    &testNatGatewayVar,
			Items: &[]ionoscloud.NatGateway{natgatewayTest.NatGateway},
		},
	}
	natgatewayProperties = resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Name:      &testNatGatewayNewVar,
			PublicIps: &[]string{testNatGatewayNewVar},
		},
	}
	natgatewayNew = resources.NatGateway{
		NatGateway: ionoscloud.NatGateway{
			Properties: &natgatewayProperties.NatGatewayProperties,
		},
	}
	testNatGatewayVar    = "test-natgateway"
	testNatGatewayNewVar = "test-new-natgateway"
	testNatGatewayErr    = errors.New("natgateway test error")
)

func TestPreRunDcIdsNatGatewayProperties(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNatGatewayVar)
		err := PreRunDcIdsNatGatewayIps(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdsNatGatewayPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcIdsNatGatewayIps(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNatGatewayIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		err := PreRunDcNatGatewayIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNatGatewayIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNatGatewayIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgCols), defaultNatGatewayCols)
		rm.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgateways, nil, nil)
		err := RunNatGatewayList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		rm.NatGateway.EXPECT().List(testNatGatewayVar).Return(natgateways, nil, testNatGatewayErr)
		err := RunNatGatewayList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, nil)
		err := RunNatGatewayGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, nil)
		rm.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, nil)
		err := RunNatGatewayGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, testNatGatewayErr)
		err := RunNatGatewayGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		rm.NatGateway.EXPECT().Get(testNatGatewayVar, testNatGatewayVar).Return(&natgatewayTestGet, nil, testNatGatewayErr)
		err := RunNatGatewayGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayVar})
		rm.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, nil, nil)
		err := RunNatGatewayCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayVar})
		rm.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, &testResponse, nil)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayVar})
		rm.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, nil, testNatGatewayErr)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().Create(testNatGatewayVar, natgatewayTest).Return(&natgatewayTest, nil, nil)
		err := RunNatGatewayCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), testNatGatewayNewVar)
		rm.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, nil, nil)
		err := RunNatGatewayUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayNewVar})
		rm.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, nil, testNatGatewayErr)
		err := RunNatGatewayUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testNatGatewayNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIps), []string{testNatGatewayNewVar})
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().Update(testNatGatewayVar, testNatGatewayVar, natgatewayProperties).Return(&natgatewayNew, nil, nil)
		err := RunNatGatewayUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		rm.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		rm.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, testNatGatewayErr)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, nil)
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNatGatewayDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.NatGateway.EXPECT().Delete(testNatGatewayVar, testNatGatewayVar).Return(nil, nil)
		err := RunNatGatewayDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNatGatewayDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testNatGatewayVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNatGatewayId), testNatGatewayVar)
		cfg.Stdin = os.Stdin
		err := RunNatGatewayDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNatGatewaysCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Name"})
	getNatGatewaysCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNatGatewaysColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("natgateway", config.ArgCols), []string{"Unknown"})
	getNatGatewaysCols(core.GetGlobalFlagName("natgateway", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetNatGatewaysIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getNatGatewaysIds(w, testNatGatewayVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
