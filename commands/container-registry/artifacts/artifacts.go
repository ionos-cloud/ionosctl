package artifacts

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "TotalVulnerabilities", "FixableVulnerabilities", "MediaType"}
	allCols     = []string{
		"Id", "Repository", "PushCount", "PullCount", "LastPushed", "TotalVulnerabilities",
		"FixableVulnerabilities", "MediaType", "URN",
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
