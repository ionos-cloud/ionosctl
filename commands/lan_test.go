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
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		err := PreRunGlobalDcId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdLanId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		err := PreRunGlobalDcIdLanId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdLanIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), "")
		err := PreRunGlobalDcIdLanId(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		rm.Lan.EXPECT().List(testLanVar).Return(ls, nil, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		rm.Lan.EXPECT().List(testLanVar).Return(ls, nil, testLanErr)
		err := RunLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, nil)
		err := RunLanGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, testLanErr)
		err := RunLanGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanVar)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, testLanErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanVar)
		rm.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testLanNewVar)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetLansCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("lan", config.ArgCols), []string{"Name"})
	getLansCols(builder.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLansColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("lan", config.ArgCols), []string{"Unknown"})
	getLansCols(builder.GetGlobalFlagName("lan", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(builder.GetGlobalFlagName("lan", config.ArgDataCenterId), testLanVar)
	getLansIds(w, testLanVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
