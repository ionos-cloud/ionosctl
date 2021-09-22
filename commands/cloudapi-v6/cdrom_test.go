package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testCdroms = resources.Cdroms{
		Cdroms: ionoscloud.Cdroms{
			Items: &[]ionoscloud.Image{testImage.Image},
		},
	}
	testCdromVar = "test-cdrom"
	testCdromErr = errors.New("cdrom test error")
)

func TestPreRunDcServerCdromIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		err := PreRunDcServerCdromIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerCdromIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcServerCdromIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, &testResponse, nil)
		err := RunServerCdromAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar).Return(testCdroms, &testResponse, nil)
		err := RunServerCdromsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar).Return(testCdroms, nil, testCdromErr)
		err := RunServerCdromsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, &testResponse, nil)
		err := RunServerCdromGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testResponseErr, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunCdromDetachAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
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
