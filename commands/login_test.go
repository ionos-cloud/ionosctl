package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunLoginUser_BufferUserErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.Name, "user"), "")
		viper.Set(builder.GetFlagName(cfg.Name, "password"), "test")
		cfg.Stdin = bytes.NewReader([]byte("test@test.com\n"))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == "401 Unauthorized")
	})
}

func TestRunLoginUser_BufferErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.Name, "user"), "")
		viper.Set(builder.GetFlagName(cfg.Name, "password"), "test")
		cfg.Stdin = bytes.NewReader([]byte("test@test.com"))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoginUser_UnauthorizedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.Name, "user"), "test@test.com")
		viper.Set(builder.GetFlagName(cfg.Name, "password"), "test")
		err := RunLoginUser(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == "401 Unauthorized")
	})
}

func TestRunLoginUser_BufferPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.Name, "user"), "test@test.com")
		viper.Set(builder.GetFlagName(cfg.Name, "password"), "")
		err := RunLoginUser(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == "inappropriate ioctl for device")
	})
}
