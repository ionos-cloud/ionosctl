package artifacts

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Repository", JSONPath: "properties.repositoryName"},
	{Name: "PushCount", JSONPath: "metadata.pushCount"},
	{Name: "PullCount", JSONPath: "metadata.pullCount"},
	{Name: "LastPushed", JSONPath: "metadata.lastPushedAt"},
	{Name: "TotalVulnerabilities", JSONPath: "metadata.vulnTotalCount", Default: true},
	{Name: "FixableVulnerabilities", JSONPath: "metadata.vulnFixableCount", Default: true},
	{Name: "MediaType", JSONPath: "properties.mediaType", Default: true},
	{Name: "URN", JSONPath: "metadata.resourceURN"},
	{Name: "RegistryId"},
}

func ArtifactsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "artifacts",
			Aliases: []string{"a", "art", "artifact"},
			Short:   "Artifacts Operations",
			Long: "Manage container registry artifacts. " +
				"Artifacts are the individual files stored in a repository.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(ArtifactsListCmd())
	cmd.AddCommand(ArtifactsGetCmd())

	return cmd
}

func ArtifactsIds(registryId string, repositoryName string) []string {
	artifacts, _, err := client.Must().RegistryClient.ArtifactsApi.RegistriesRepositoriesArtifactsGet(
		context.Background(), registryId, repositoryName,
	).Execute()
	if err != nil {
		return nil
	}

	t := table.New(allCols, table.WithPrefix("items"))
	if err := t.Extract(artifacts); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Id").AddInfo("MediaType").ToString()
}
