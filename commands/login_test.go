package commands

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "test@ionos.com"
	testPassword = "testPwd"
	testToken    = "testToken"
)

func TestPreRunLoginCmd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), testUsername)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			err := PreRunLoginCmd(cfg)
			assert.NoError(t, err)
		},
	)
}

func TestPreRunLoginCmdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(
		t, w, func(cfg *core.PreCommandConfig) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), testUsername)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgToken), testToken)
			err := PreRunLoginCmd(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserTokenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgToken), testToken)
			err := RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserBufferUserErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), "")
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			cfg.Stdin = bytes.NewReader([]byte(testUsername + "\n"))
			err := RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserBufferErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), "")
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			cfg.Stdin = bytes.NewReader([]byte(testUsername))
			err := RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserUnauthorizedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), testUsername)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			err := RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserBufferPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), testUsername)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), "")
			err := RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}

func TestRunLoginUserConfigSet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(
		t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
			viper.Reset()
			err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
			assert.NoError(t, err)
			err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
			assert.NoError(t, err)
			err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
			assert.NoError(t, err)
			viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), testUsername)
			viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
			err = RunLoginUser(cfg)
			assert.Error(t, err)
		},
	)
}
