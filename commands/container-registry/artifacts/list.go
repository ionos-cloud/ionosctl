package artifacts

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ArtifactsListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "artifacts",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List registry or repository artifacts",
			LongDesc:   "List all artifacts in a registry or repository",
			Example:    "ionosctl container-registry artifacts list",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag("repository", "", "", "Name of the repository to list artifacts from")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all artifacts in the registry")

	return cmd
}

func PreCmdList(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{"registry-id", "repository"},
		[]string{"registry-id", constants.ArgAll},
	)
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	regId := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
	defCols := defaultCols

	var arts interface{}
	var err error
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		fmt.Println("Listing all artifacts in registry")
		arts, _, err = client.Must().RegistryClient.ArtifactsApi.RegistriesArtifactsGet(c.Context, regId).Execute()
		if err != nil {
			return err
		}

		defCols = append(defCols, "Repository")
	} else {
		repo := viper.GetString(core.GetFlagName(c.NS, "repository"))

		arts, _, err = client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsGet(
			c.Context, regId, repo,
		).Execute()
		if err != nil {
			return err
		}
	}

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.ContainerRegistryArtifact, arts,
		tabheaders.GetHeaders(allCols, defCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
