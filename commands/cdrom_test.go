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
	testCdroms = v5.Cdroms{
		Cdroms: ionoscloud.Cdroms{
			Items: &[]ionoscloud.Image{testImage.Image},
		},
	}
	testCdromVar = "test-cdrom"
	testCdromErr = errors.New("cdrom test error")
)

func TestRunServerCdromAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, &testResponse, nil)
		err := RunServerCdromAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, nil, nil)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		rm.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar).Return(testCdroms, &testResponse, nil)
		err := RunServerCdromsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		rm.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar).Return(testCdroms, nil, testCdromErr)
		err := RunServerCdromsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		rm.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, &testResponse, nil)
		err := RunServerCdromGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		rm.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testResponseErr, nil)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunCdromDetachAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestGetAttachedCdromsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getAttachedCdromsIds(w, testVolumeVar, testVolumeVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetImagesCdromIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getImagesCdromIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
