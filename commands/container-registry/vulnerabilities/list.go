package vulnerabilities

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/artifacts"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/repository"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, "repository", constants.FlagArtifactId,
	); err != nil {
		return err
	}

	return nil
}

func CmdList(c *core.CommandConfig) error {
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

	return c.Out(table.Sprint(allCols, vulnerabilities, cols, table.WithPrefix("items")))
}
