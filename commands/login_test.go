//go:build integration
// +build integration

package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "test@ionos.com"
	testPassword = "testPwd"
	testToken    = "testToken"
)

func TestRunLoginUserBufferErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgUser), "")
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testPassword)
		cfg.Stdin = bytes.NewReader([]byte(testUsername))
		err := RunLoginUser(cfg)
		assert.Error(t, err)
	})
}
