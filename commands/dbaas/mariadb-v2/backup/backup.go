package backup

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/backup/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var backupCols = []table.Column{
	{Name: "BackupId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "ClusterName", JSONPath: "properties.clusterName", Default: true},
	{Name: "Version", JSONPath: "properties.mariadbClusterVersion", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "EarliestRecoveryTargetTime", JSONPath: "properties.earliestRecoveryTargetTime", Default: true},
}

func BackupCmd() *core.Command {
	backupCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "backup",
			Aliases: []string{"b", "backups"},
			Short:   "MariaDB Backup Operations",
			Long: `The sub-commands of ` + "`ionosctl dbaas mariadb-v2 backup`" + ` allow you to view MariaDB Cluster Backups.

A backup is not a single point in time — it represents a recovery WINDOW starting at its earliestRecoveryTargetTime and extending to the present. The cluster can be restored to any moment within that window. Use ` + "`--recovery-time`" + ` on ` + "`cluster restore`" + ` (or ` + "`cluster create --backup-id`" + `) to zoom into a specific point; omit it to use the latest point.`,
			TraverseChildren: true,
		},
	}

	backupCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(backupCols))
	_ = backupCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(backupCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	backupCmd.AddCommand(BackupListCmd())
	backupCmd.AddCommand(BackupGetCmd())
	backupCmd.AddCommand(location.LocationCmd())

	return backupCmd
}
