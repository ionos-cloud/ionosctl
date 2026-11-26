package version

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var versionCols = []table.Column{
	{Name: "VersionId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.version", Default: true},
	{Name: "Status", JSONPath: "properties.status", Default: true},
	{Name: "Comment", JSONPath: "properties.comment", Default: true},
	{Name: "CanUpgradeTo", JSONPath: "properties.canUpgradeTo", Default: true},
	{Name: "CreationFreezeDate", Format: optionalDate("properties.creationFreezeDate"), Default: true},
	{Name: "DeprecationDate", Format: optionalDate("properties.deprecationDate"), Default: true},
}

// zeroDate is what the API returns for a version date that is not set.
const zeroDate = "0001-01-01T00:00:00Z"

func optionalDate(path string) table.FormatFunc {
	return func(item map[string]any) any {
		if v := table.Navigate(item, path); v != nil && v != zeroDate {
			return v
		}
		return ""
	}
}

func VersionCmd() *core.Command {
	versionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "MongoDB Version Operations",
			Long:             "The sub-commands of `ionosctl dbaas mongo-v2 version` allow you to list available MongoDB Versions.",
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
