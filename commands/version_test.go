package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/stretchr/testify/assert"
)

const (
	testLatestReleaseUrl    = "https://github.com/ionos-cloud/ionosctl/releases/latest"
	testLatestReleaseApiUrl = "https://api.github.com/ionos-cloud/ionosctl/releases/latest"
)

func TestRunVersion(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		err := RunVersion(cfg)
		assert.NoError(t, err)
	})
}

func TestGetGithubLatestRelease(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	_, err := getGithubLatestVersion(latestReleaseUrl)
	assert.NoError(t, err)
}

func TestGetGithubLatestReleaseErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	_, err := getGithubLatestVersion("")
	assert.Error(t, err)
}

func TestGetGithubLatestReleaseJsonErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	_, err := getGithubLatestVersion(testLatestReleaseUrl)
	assert.Error(t, err)
}

func TestGetGithubLatestReleaseTagErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	clierror.ErrAction = func() { return }
	_, err := getGithubLatestVersion(testLatestReleaseApiUrl)
	assert.Error(t, err)
}
