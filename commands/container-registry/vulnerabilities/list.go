package vulnerabilities

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
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

	c.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag("repository", "", "", "Name of the repository to retrieve artifact from")
	_ = c.Command.RegisterFlagCompletionFunc(
		"repository", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return repository.RepositoryNames(viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagArtifactId, "", "", "ID/digest of the artifact")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagArtifactId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return artifacts.ArtifactsIds(
					viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId)),
					viper.GetString(core.GetFlagName(c.NS, "repository")),
				),
				cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddSetFlag(
		cloudapiv6.ArgOrderBy, "", "-score", []string{
			"-score", "-severity", "-publishedAt", "-updatedAt", "-fixable", "score",
			"severity", "publishedAt", "updatedAt", "fixable",
		}, cloudapiv6.ArgOrderByDescription,
	)
	c.AddStringSliceFlag(
		cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription,
	)

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, "repository", constants.FlagArtifactId,
	); err != nil {
		return err
	}

	return query.ValidateFilters(c, []string{"severity", "fixable"}, "Filters available: severity, fixable")
}

func CmdList(c *core.CommandConfig) error {
	// TODO alex: verify we can still filter by "severity" and "fixable"

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	registryId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	repository := viper.GetString(core.GetFlagName(c.NS, "repository"))
	artifactId := viper.GetString(core.GetFlagName(c.NS, constants.FlagArtifactId))

	vulnerabilities, _, err := client.Must().RegistryClient.ArtifactsApi.
		RegistriesRepositoriesArtifactsVulnerabilitiesGet(
			context.Background(), registryId, repository, artifactId).Execute()
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
