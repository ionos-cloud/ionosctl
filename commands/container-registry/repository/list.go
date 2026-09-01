package repository

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepositoryListCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "repository",
			Verb:       "list",
			Aliases:    []string{"ls", "l"},
			ShortDesc:  "List all repositories in a registry",
			LongDesc:   "List all repositories in a registry, with per-repository usage stats (artifact count, pull/push counts) and the highest vulnerability severity seen among each repository's artifacts. Repositories are created implicitly the first time an image is pushed.",
			Example:    "ionosctl container-registry repository list --registry-id REGISTRY_ID",
			PreCmdRun:  PreCmdList,
			CmdRun:     CmdList,
			InitClient: true,
		},
	)

	c.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry whose repositories to list")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func PreCmdList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId); err != nil {
		return err
	}

	return nil
}

func CmdList(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))

	repos, _, err := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesGet(
		context.Background(), regId).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(repos)
}
