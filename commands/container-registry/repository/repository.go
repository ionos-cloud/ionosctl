package repository

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
)

var (
	note = "NOTE: This command's behavior will be replaced by `ionosctl container-registry repository delete` in the" +
		" future. Please use that command instead.\n"

	allCols = []table.Column{
		{Name: "Id", JSONPath: "id", Default: true},
		{Name: "Name", JSONPath: "properties.name", Default: true},
		{Name: "LastSeverity", JSONPath: "metadata.lastSeverity", Default: true},
		{Name: "ArtifactCount", JSONPath: "metadata.artifactCount", Default: true},
		{Name: "PullCount", JSONPath: "metadata.pullCount", Default: true},
		{Name: "PushCount", JSONPath: "metadata.pushCount", Default: true},
		{Name: "LastPushedAt", JSONPath: "metadata.lastPushedAt"},
		{Name: "LastPulledAt", JSONPath: "metadata.lastPulledAt"},
		{Name: "URN", JSONPath: "metadata.resourceURN"},
	}
)

func RegRepoDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "repository",
			Verb:      "repository",
			Aliases:   []string{"rd", "del", "repo", "rep-del", "repository-delete"},
			ShortDesc: note + "Delete all contents of a repository (deprecated form).",
			LongDesc: note + `Delete all contents (all artifacts and tags) of a repository, effectively removing the repository.

The registry V2 API only allows manifests and blobs to be deleted individually, not whole repositories, so this command performs that cleanup for you as a convenience. This is irreversible.`,
			Example: "ionosctl container-registry repository-delete --registry-id REGISTRY_ID --name REPOSITORY_NAME",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				fmt.Fprint(c.Command.Command.ErrOrStderr(), note)
				return PreCmdDelete(c)
			},
			CmdRun: func(c *core.CommandConfig) error {
				fmt.Fprint(c.Command.Command.ErrOrStderr(), note)
				return CmdDelete(c)
			},
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

	cmd.AddColsFlag(allCols)

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
		repos.Items, func(repo containerregistry.RepositoryRead) string {
			return repo.Properties.Name
		},
	)
}
