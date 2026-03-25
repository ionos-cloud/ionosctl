package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var versionCols = []table.Column{
	{Name: "Name", JSONPath: "properties.version", Default: true},
	{Name: "Status", JSONPath: "properties.status", Default: true},
	{Name: "Comment", JSONPath: "properties.comment"},
	{Name: "CanUpgradeTo", JSONPath: "properties.canUpgradeTo"},
}

var allVersionCols = table.AllCols(versionCols)
var defaultVersionCols = table.DefaultCols(versionCols)

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
