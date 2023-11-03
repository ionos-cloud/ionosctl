package commands

import (
	"bufio"
	"bytes"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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

func TestGetGithubLatestRelease(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expected  string
		shouldErr bool
	}{
		{
			name:      "Empty URL",
			url:       "",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Non-API URL",
			url:       "https://github.com/user/repo/releases/latest",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Invalid Tag",
			url:       "https://api.github.com/user/repo/releases/latest",
			expected:  "",
			shouldErr: true,
		},
		{
			name: "Valid Tag",
			url:  latestGhApiReleaseUrl,
			// Regex pattern for semver
			expected:  `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latest, err := getGithubLatestRelease(tt.url)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				match, _ := regexp.MatchString(tt.expected, latest)
				assert.True(t, match, "Expected a valid semver version")
			}
		})
	}
}
