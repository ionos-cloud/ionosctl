package artifacts

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ArtifactsGetCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "artifacts",
			Verb:       "get",
			ShortDesc:  "Retrieve an artifacts",
			LongDesc:   "Retrieve an artifact from a repository",
			Example:    "ionosctl container-registry artifacts get",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	c.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = c.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag("repository", "", "", "Name of the repository to retrieve artifact from")
	_ = c.Command.RegisterFlagCompletionFunc(
		"repository", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return repository.RepositoryNames(viper.GetString(core.GetFlagName(c.NS, "registry-id"))),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag("artifact-id", "", "", "ID/digest of the artifact")
	_ = c.Command.RegisterFlagCompletionFunc(
		"artifact-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ArtifactsIds(
					viper.GetString(core.GetFlagName(c.NS, "registry-id")),
					viper.GetString(core.GetFlagName(c.NS, "repository")),
				),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, "registry-id", "repository", "artifact-id")
}

func CmdGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	regId := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
	repo := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artId := viper.GetString(core.GetFlagName(c.NS, "artifact-id"))

	var arts interface{}
	var err error

	arts, _, err = client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsFindByDigest(
		c.Context, regId, repo, artId,
	).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.ContainerRegistryArtifact, arts,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
