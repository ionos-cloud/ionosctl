package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
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
			ShortDesc: "Delete all repository contents.",
			LongDesc: "Delete all repository contents.\n" +
				"The registry V2 API allows manifests and blobs to be deleted individually, but " +
				"it is not possible to remove an entire repository. This operation is provided for convenience",
			Example:    "ionosctl container-registry delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the repository to delete")
	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "Registry ID")
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Repository is being deleted"))

	return nil
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
