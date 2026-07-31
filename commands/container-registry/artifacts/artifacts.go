package artifacts

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
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
}

func ArtifactsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "artifacts",
			Aliases: []string{"a", "art", "artifact"},
			Short:   "Inspect artifacts stored in a registry",
			Long: `An artifact is a single object stored in a repository - a Docker image, an image manifest, or an OCI artifact - addressed by its content digest (sha256:...). One repository holds many artifacts; tags such as 'latest' are pointers to a specific artifact digest.

These commands are read-only: artifacts are created by 'docker push' and removed by garbage collection or by deleting the repository. Each artifact carries usage stats (pull/push counts) and, when vulnerability scanning is enabled on the registry, a count of total and fixable vulnerabilities (drill in with 'container-registry vulnerabilities list').`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

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
