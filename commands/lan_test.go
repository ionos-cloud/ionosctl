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
	lp           = ionoscloud.LanPost{
		Id: &testLanVar,
		Properties: &ionoscloud.LanPropertiesPost{
			Name:       &testLanVar,
			IpFailover: nil,
			Pcc:        &testLanVar,
			Public:     &publicLan,
		},
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
				Pcc:        &testLanVar,
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

func TestPreRunGlobalDcIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		err := PreRunGlobalDcIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())
	})
}

func TestPreRunGlobalDcIdLanIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		err := PreRunGlobalDcIdLanIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdLanIdValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdLanIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), "")
		err = PreRunGlobalDcIdLanIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgLanId).Error())
	})
}

func TestRunLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		rm.Lan.EXPECT().List(testLanVar).Return(ls, nil, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{l}, nil, nil)
		err := RunLanGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		rm.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{l}, nil, testLanErr)
		err := RunLanGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		rm.Lan.EXPECT().Create(testLanVar, testLanVar, publicLan).Return(&resources.LanPost{lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanCreate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		rm.Lan.EXPECT().Create(testLanVar, testLanVar, publicLan).Return(&resources.LanPost{lp}, nil, testLanErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicLan)
		rm.Lan.EXPECT().Create(testLanVar, testLanVar, publicLan).Return(&resources.LanPost{lp}, nil, nil)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanUpdate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgServerId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanName), testLanNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanPublic), publicNewLan)
		rm.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, nil)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDelete_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLanId), testLanVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDelete_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
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
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("lan", config.ArgCols), []string{"Name"})
	getLansCols(builder.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLansCols_Err(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

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
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(builder.GetGlobalFlagName("lan", config.ArgDataCenterId), testLanVar)
	getLansIds(w, "lan")
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
