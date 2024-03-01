package version

import (
	"fmt"
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
		versionParts := strings.Split(info.Main.Version, "-")
		if len(versionParts) >= 3 {
			commitHash := versionParts[2]
			if len(commitHash) > 7 {
				commitHash = commitHash[:7]
			}
			return fmt.Sprintf("DEV-%s", commitHash)
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
