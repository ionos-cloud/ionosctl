package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var versionCols = []table.Column{
	{Name: "Name", JSONPath: "properties.version", Default: true},
	{Name: "Status", JSONPath: "properties.status", Default: true},
	{Name: "Comment", JSONPath: "properties.comment", Default: true},
	{Name: "CanUpgradeTo", JSONPath: "properties.canUpgradeTo", Default: true},
}

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

	versionCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(versionCols))
	_ = versionCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(versionCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	versionCmd.AddCommand(VersionListCmd())
	versionCmd.AddCommand(VersionGetCmd())

	return versionCmd
}
