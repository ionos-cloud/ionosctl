package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/core"
)

const (
	latestGhApiReleaseUrl = "https://api.github.com/repos/ionos-cloud/ionosctl/releases/latest"
	latestGhReleaseUrl    = "https://github.com/ionos-cloud/ionosctl/releases/latest"
)

func version() *core.Command {
	ctx := context.TODO()
	versionCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "version",
		Resource:   "version",
		Verb:       "version",
		ShortDesc:  "Show the current version",
		LongDesc:   "The `ionosctl version` command displays the current version of the ionosctl software and the latest Github release.",
		Example:    "ionosctl version",
		PreCmdRun:  noPreRun,
		CmdRun:     RunVersion,
		InitClient: false,
	})

	return versionCmd
}

func RunVersion(c *core.CommandConfig) error {
	err := c.Printer.Print("You are currently using ionosctl version: " + rootCmd.Command.Version)
	if err != nil {
		return err
	}

	// Get and Print the latest Github Release for IONOS Cloud CLI
	latestGhRelease, err := getGithubLatestRelease(latestGhApiReleaseUrl)
	if err != nil {
		return err
	}
	return c.Printer.Print(fmt.Sprintf("Latest ionosctl Github release: %s\nFor more information, go to: %s", latestGhRelease, latestGhReleaseUrl))
}

// getGithubLatestRelease retrieves the latest official release
// from Github or returns an error, if any.
func getGithubLatestRelease(u string) (string, error) {
	res, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var m map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&m); err != nil {
		return "", err
	}

	tn, ok := m["tag_name"]
	if !ok {
		return "", errors.New("error finding tag_name in response")
	}
	v := strings.TrimPrefix(tn.(string), "v")
	return v, nil
}
