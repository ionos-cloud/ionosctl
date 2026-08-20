package location

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var locationCols = []table.Column{
	{Name: "BackupLocationId", JSONPath: "id", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
}

func LocationCmd() *core.Command {
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc", "locations"},
			Short:            "MariaDB Backup Location Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb-v2 backup location` allow you to list the supported Object Storage locations where backups can be stored.",
			TraverseChildren: true,
		},
	}

	locationCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(locationCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(locationCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	locationCmd.AddCommand(LocationListCmd())
	locationCmd.AddCommand(LocationGetCmd())

	return locationCmd
}
