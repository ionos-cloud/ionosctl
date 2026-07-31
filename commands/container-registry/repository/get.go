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

func RepositoryGetCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "repository",
			Verb:       "get",
			ShortDesc:  "Get a repository's properties and stats",
			LongDesc:   "Get a single repository from a registry by name, including usage stats (artifact count, pull/push counts, last pushed/pulled timestamps) and the highest vulnerability severity seen among its artifacts.",
			Example:    "ionosctl container-registry repository get --registry-id REGISTRY_ID --name REPOSITORY_NAME",
			PreCmdRun:  PreCmdGet,
			CmdRun:     CmdGet,
			InitClient: true,
		},
	)

	c.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry the repository belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the repository to get (the path in <hostname>/<repository>)")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagName,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RepositoryNames(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId)),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func PreCmdGet(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId, constants.FlagName)
}

func CmdGet(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

	repo, _, err := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesFindByName(
		context.
			Background(), regId, name,
	).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(repo)
}
