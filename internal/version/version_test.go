package version

import (
	"runtime/debug"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name            string
		setup           func()
		fetcher         BuildInfoFetcher
		expectedVersion string
	}{
		{
			name: "ldflags version and release label",
			setup: func() {
				Version = "1.0.0"
				GitCommit = "abc123"
				Label = "release"
			},
			expectedVersion: "v1.0.0",
		},
		{
			name: "ldflags GitCommit only",
			setup: func() {
				Version = ""
				GitCommit = "def456"
				Label = ""
			},
			expectedVersion: "DEV-def456",
		},
		{
			name: "ldflags GitCommit modified files",
			setup: func() {
				Version = ""
				GitCommit = "def456+"
				Label = ""
			},
			expectedVersion: "DEV-def456+",
		},
		{
			name: "semantic version from build info",
			setup: func() {
				Version = ""
				GitCommit = ""
				Label = ""
			},
			fetcher: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{
						Version: "v2.0.0",
					},
				}, true
			},
			expectedVersion: "v2.0.0",
		},
		{
			name: "commit hash from build info",
			setup: func() {
				Version = ""
				GitCommit = ""
				Label = ""
			},
			fetcher: func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{
						Version: "v0.0.0-20210101000000-abcdef123456",
					},
				}, true
			},
			expectedVersion: "DEV-abcdef1",
		},
		{
			name: "unknown version",
			setup: func() {
				Version = ""
				GitCommit = ""
				Label = ""
			},
			fetcher: func() (*debug.BuildInfo, bool) {
				return nil, false
			},
			expectedVersion: "unknown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer func() {
				Version = ""
				GitCommit = ""
				Label = ""
			}()

			var got string
			if tc.fetcher != nil {
				got = Get(WithFetcher(tc.fetcher))
			} else {
				got = Get()
			}

			if got != tc.expectedVersion {
				t.Errorf("expected version %q, got %q", tc.expectedVersion, got)
			}
		})
	}
}
