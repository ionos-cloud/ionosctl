package vulnerabilities

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VulnerabilitiesListCmd() *core.Command {
	c := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "vulnerabilities",
			Verb:       "list",
			ShortDesc:  "Retrieve vulnerabilities",
			LongDesc:   "Retrieve all vulnerabilities from an artifact",
			Example:    "ionosctl container-registry vulnerabilities list",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
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
			return artifacts.ArtifactsIds(
					viper.GetString(core.GetFlagName(c.NS, "registry-id")),
					viper.GetString(core.GetFlagName(c.NS, "repository")),
				),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, "registry-id", "repository", "artifact-id")
}

func CmdList(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	registryId := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
	repository := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artifactId := viper.GetString(core.GetFlagName(c.NS, "artifact-id"))

	vulnerabilities, _, err := client.Must().RegistryClient.ArtifactsApi.
		RegistriesRepositoriesArtifactsVulnerabilitiesGet(
			context.Background(), registryId, repository, artifactId,
		).Execute()
	if err != nil {
		return err
	}

	vulnerabilitiesConverted, err := resource2table.ConvertContainerRegistryVulnerabilitiesToTable(
		vulnerabilities,
	)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		vulnerabilities, vulnerabilitiesConverted,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
