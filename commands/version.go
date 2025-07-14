package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/sdk-go-bundle/products/auth/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	sdkcertmanager "github.com/ionos-cloud/sdk-go-cert-manager"
	sdkpostgres "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/viper"
)

const (
	latestGhApiReleaseUrl = "https://api.github.com/repos/ionos-cloud/ionosctl/releases/latest"
	latestGhReleaseUrl    = "https://github.com/ionos-cloud/ionosctl/releases/latest"
)

func VersionCmd() *core.Command {
	ctx := context.TODO()
	versionCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "version",
		Resource:   "version",
		Verb:       "version",
		ShortDesc:  "Show the current version",
		LongDesc:   "The `ionosctl version` command displays the current version of the ionosctl software and the latest Github release.",
		Example:    "ionosctl version",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunVersion,
		InitClient: false,
	})
	versionCmd.AddBoolFlag(constants.ArgUpdates, "", false, "Check for latest updates for CLI")

	return versionCmd
}

func RunVersion(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(rootCmd.Command.Version))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("sdk-go "+compute.Version))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("sdk-go-dbaas-postgres "+sdkpostgres.Version))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("sdk-go-auth "+auth.Version))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("sdk-go-cert-manager "+sdkcertmanager.Version))

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgUpdates)) {
		/*
			Latest Github Release for IONOS Cloud CLI
		*/
		latestGhRelease, err := getGithubLatestRelease(latestGhApiReleaseUrl)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(
			"Latest ionosctl Github release: %s\nFor more information, go to: %s", latestGhRelease, latestGhReleaseUrl))

		return nil
	}
	return nil
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
