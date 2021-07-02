package commands

import (
	"bufio"
	"bytes"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "test@ionos.com"
	testPassword = "test"
)

func TestRunLoginUserBufferUserErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUser), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testPassword)
		cfg.Stdin = bytes.NewReader([]byte(testUsername + "\n"))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserBufferErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUser), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testPassword)
		cfg.Stdin = bytes.NewReader([]byte(testUsername))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserUnauthorizedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUser), testUsername)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testPassword)
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserBufferPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUser), testUsername)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), "")
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserConfigSet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
		assert.NoError(t, err)
		err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
		assert.NoError(t, err)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUser), testUsername)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testPassword)
		err = RunLoginUser(cfg)
		assert.Error(t, err)
	})
}
