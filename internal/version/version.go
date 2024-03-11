package version

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

// Variables sent at link time via ldflags
var (
	Version   string
	GitCommit string
	Label     string
)

// Get returns the current CLI version, preferring the value set via ldflags
func Get() string {
	if Version != "" || GitCommit != "" {
		return getVersionViaLdFlags()
	}

	info, ok := debug.ReadBuildInfo()
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
