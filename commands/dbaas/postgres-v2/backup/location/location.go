package location

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var backupLocationCols = []table.Column{
	{Name: "LocationId", JSONPath: "id", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
}

func BackupLocationCmd() *core.Command {
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc"},
			Short:            "PostgreSQL Backup Location Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres-v2 backup location` allow you to list and get DBaaS PostgreSQL Backup Locations.",
			TraverseChildren: true,
		},
	}

	locationCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(backupLocationCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(backupLocationCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	locationCmd.AddCommand(BackupLocationListCmd())
	locationCmd.AddCommand(BackupLocationGetCmd())

	return locationCmd
}
