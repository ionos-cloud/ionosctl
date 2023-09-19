package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunVersion(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgUpdates), false)
		err := RunVersion(cfg)
		assert.NoError(t, err)
	})
}

func TestGetGithubLatestReleaseErr(t *testing.T) {
	defer func(a func()) { die.ErrAction = a }(die.ErrAction)
	die.ErrAction = func() { return }
	_, err := getGithubLatestRelease("")
	assert.Error(t, err)
}

func TestGetGithubLatestReleaseJsonErr(t *testing.T) {
	defer func(a func()) { die.ErrAction = a }(die.ErrAction)
	die.ErrAction = func() { return }
	_, err := getGithubLatestRelease(latestGhReleaseUrl)
	assert.Error(t, err)
}

func TestGetGithubLatestReleaseTagErr(t *testing.T) {
	defer func(a func()) { die.ErrAction = a }(die.ErrAction)
	die.ErrAction = func() { return }
	_, err := getGithubLatestRelease("https://api.github.com/ionos-cloud/ionosctl/releases/latest")
	assert.Error(t, err)
}
