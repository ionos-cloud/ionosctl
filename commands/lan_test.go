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
	publicLan    = true
	publicNewLan = false
	lanPostTest  = ionoscloud.LanPost{
		Properties: &ionoscloud.LanPropertiesPost{
			Name:       &testLanVar,
			IpFailover: nil,
			Pcc:        &testLanVar,
			Public:     &publicLan,
		},
	}
	lp = ionoscloud.LanPost{
		Id:         &testLanVar,
		Properties: lanPostTest.Properties,
		Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
	}
	l = ionoscloud.Lan{
		Id: &testLanVar,
		Properties: &ionoscloud.LanProperties{
			Name: &testLanVar,
			Pcc:  &testLanVar,
		},
	}
	lanProperties = resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			Name:   &testLanNewVar,
			Pcc:    &testLanNewVar,
			Public: &publicNewLan,
		},
	}
	lanNew = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testLanVar,
			Properties: &ionoscloud.LanProperties{
				Name:       lanProperties.LanProperties.Name,
				Public:     lanProperties.LanProperties.Public,
				IpFailover: nil,
				Pcc:        &testLanNewVar,
			},
		},
	}
	ls = resources.Lans{
		Lans: ionoscloud.Lans{
			Id:    &testLanVar,
			Items: &[]ionoscloud.Lan{l},
		},
	}
	testLanVar    = "test-lan"
	testLanNewVar = "test-new-lan"
	testLanErr    = errors.New("lan test: error occurred")
)

func TestPreRunGlobalDcId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		err := PreRunGlobalDcId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdLanId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		err := PreRunGlobalDcIdLanId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdLanIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdLanId(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		rm.Lan.EXPECT().List(testLanVar).Return(ls, nil, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		rm.Lan.EXPECT().List(testLanVar).Return(ls, nil, testLanErr)
		err := RunLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, nil)
		err := RunLanGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, testLanErr)
		err := RunLanGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicLan)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanVar)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, testLanErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanVar)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testResponse, nil)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetLansCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"Name"})
	getLansCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLansColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"Unknown"})
	getLansCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetLansIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgDataCenterId), testLanVar)
	getLansIds(w, testLanVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
