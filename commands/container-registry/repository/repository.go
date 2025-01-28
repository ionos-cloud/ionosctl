package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var (
	note = "NOTE: This command's behavior will be replaced by `ionosctl container-registry repository delete` in the" +
		" future. Please use that command instead.\n"

	defaultCols = []string{"Id", "Name", "LastSeverity", "ArtifactCount", "PullCount", "PushCount"}
	allCols     = []string{
		"Id", "Name", "LastSeverity", "ArtifactCount", "PullCount", "PushCount", "LastPushedAt",
		"LastPulledAt", "URN",
	}
)

func RegRepoDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "repository",
			Verb:      "repository",
			Aliases:   []string{"rd", "del", "repo", "rep-del", "repository-delete"},
			ShortDesc: note + "Delete all repository contents.",
			LongDesc: note + "Delete all repository contents. " +
				"The registry V2 API allows manifests and blobs to be deleted " +
				"individually but it is not possible to remove an entire repository. This operation is provided for " +
				"convenience",
			Example: "ionosctl container-registry repository-delete --registry-id [REGISTRY-ID], --name [REPOSITORY-NAME]",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), note)
				return PreCmdDelete(c)
			},
			CmdRun: func(c *core.CommandConfig) error {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), note)
				return CmdDelete(c)
			},
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

	cmd.AddCommand(RepositoryDeleteCmd())
	cmd.AddCommand(RepositoryListCmd())
	cmd.AddCommand(RepositoryGetCmd())

	return cmd
}

func RepositoryNames(registryId string) []string {
	repos, _, err := client.Must().RegistryClient.RepositoriesApi.RegistriesRepositoriesGet(
		context.Background(),
		registryId,
	).Execute()
	if err != nil {
		return nil
	}

	return functional.Map(
		*repos.Items, func(repo containerregistry.RepositoryRead) string {
			return *repo.Properties.Name
		},
	)
}
