package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultVersionCols = []string{"Name", "Status"}
	allVersionCols     = []string{"Name", "Status", "Comment", "CanUpgradeTo"}
)

func VersionCmd() *core.Command {
	versionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "PostgreSQL Version Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres version` allow you to list available PostgreSQL Versions.",
			TraverseChildren: true,
		},
	}

	versionCmd.AddCommand(VersionListCmd())
	versionCmd.AddCommand(VersionGetCmd())

	return versionCmd
}
