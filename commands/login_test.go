package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "user"), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "password"), testPassword)
		cfg.Stdin = bytes.NewReader([]byte(testUsername + "\n"))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == "401 Unauthorized")
	})
}

func TestRunLoginUserBufferErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "user"), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "password"), testPassword)
		cfg.Stdin = bytes.NewReader([]byte(testUsername))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserUnauthorizedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "user"), testUsername)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "password"), testPassword)
		err := RunLoginUser(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == "401 Unauthorized")
	})
}

func TestRunLoginUserBufferPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "user"), testUsername)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "password"), "")
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUserConfigSet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "user"), testUsername)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, "password"), testPassword)
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}
