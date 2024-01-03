package vulnerabilities

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "DataSource", "Score", "Severity", "Fixable", "PublishedAt"}
	allCols     = []string{
		"Id", "DataSource", "Score", "Severity", "Fixable", "PublishedAt", "UpdatedAt", "Affects", "Description",
		"Recommendations", "References", "Href",
	}
)

func VulnerabilitiesCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vulnerabilities",
			Aliases:          []string{"v", "vuln", "vulnerability"},
			Short:            "Vulnerabilities Operations",
			Long:             "Manage container registry vulnerabilities.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(VulnerabilitiesGetCmd())
	cmd.AddCommand(VulnerabilitiesListCmd())

	return cmd
}
