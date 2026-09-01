package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepositoryDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "repository",
			Verb:      "delete",
			Aliases:   []string{"d", "del"},
			ShortDesc: "Delete all contents of a repository",
			LongDesc: `Delete all contents (every artifact and tag) of a repository, effectively removing the repository from the registry.

The registry V2 API only allows manifests and blobs to be deleted individually, not whole repositories, so this command performs that cleanup for you as a convenience. This is irreversible.`,
			Example:    "ionosctl container-registry repository delete --registry-id REGISTRY_ID --name REPOSITORY_NAME",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the repository to delete (the path in <hostname>/<repository>)")
	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry the repository belongs to")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	repoName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

	_, err := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesDelete(context.Background(), regId, repoName).Execute()
	if err != nil {
		return fmt.Errorf("failed deleting repository %s: %w", repoName, err)
	}

	c.Msg("Repository is being deleted")

	return nil
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
