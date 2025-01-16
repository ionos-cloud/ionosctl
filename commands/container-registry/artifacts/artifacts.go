package artifacts

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "TotalVulnerabilities", "FixableVulnerabilities", "MediaType"}
	allCols     = []string{
		"Id", "Repository", "PushCount", "PullCount", "LastPushed", "TotalVulnerabilities",
		"FixableVulnerabilities", "MediaType", "URN", "RegistryId",
	}
)

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

	artifactsConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.ContainerRegistryArtifact, artifacts)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(artifactsConverted, "Id").AddInfo("MediaType").ToString()
}
