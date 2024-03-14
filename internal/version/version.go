package version

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

type Option func(Config) Config
type Config struct {
	ReadBuildInfo BuildInfoFetcher
}

// BuildInfoFetcher specifies a custom build info fetcher, only really exists for testing
type BuildInfoFetcher func() (*debug.BuildInfo, bool)

// Variables sent at link time via ldflags
var (
	Version   string
	GitCommit string
	Label     string
)

// Get returns the current CLI version, preferring the value set via ldflags
// Examples:
//   - v1.0.0
//   - DEV-abcdef1
//   - DEV-abcdef1+ (if installed via `make install` with uncommitted changes)
//   - unknown
func Get(options ...Option) string {
	cfg := Config{
		ReadBuildInfo: debug.ReadBuildInfo, // Default fetcher
	}
	for _, option := range options {
		cfg = option(cfg)
	}

	if Version != "" || GitCommit != "" {
		return getVersionViaLdFlags()
	}

	info, ok := cfg.ReadBuildInfo()
	if ok {
		// If installed via a known tag using `go install`, return the version directly
		if isSemanticVersion(info.Main.Version) {
			return "v" + strings.TrimLeft(info.Main.Version, "v ")
		}

		// Installed via `go install` with a commit hash, return a dev version
		// Example: v0.0.0-20210101000000-abcdef123456 -> DEV-abcdef1
		versionParts := strings.Split(info.Main.Version, "-")
		if len(versionParts) >= 3 {
			versionOrHash := versionParts[2]
			if isSemanticVersion(versionOrHash) {
				return "v" + strings.TrimLeft(versionOrHash, "v ")
			}
			if len(versionOrHash) > 7 {
				versionOrHash = versionOrHash[:7]
			}
			return fmt.Sprintf("DEV-%s", versionOrHash)
		}
	}
	return "unknown"
}

func getVersionViaLdFlags() string {
	if Label == "release" {
		return "v" + strings.TrimLeft(Version, "v ")
	}
	return fmt.Sprintf("DEV-%s", GitCommit)
}

func isSemanticVersion(version string) bool {
	semanticVersionPattern := `^v\d+\.\d+\.\d+$`
	matched, _ := regexp.MatchString(semanticVersionPattern, version)
	return matched
}

// WithFetcher allows setting a custom fetcher for build info
func WithFetcher(fetcher BuildInfoFetcher) Option {
	return func(c Config) Config {
		if fetcher != nil {
			c.ReadBuildInfo = fetcher
		}
		return c
	}
}
