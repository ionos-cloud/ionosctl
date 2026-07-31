package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/snapshot/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var snapshotCols = []table.Column{
	{Name: "SnapshotId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "ClusterName", JSONPath: "properties.clusterName", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.datacenterId", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "ClusterVersion", JSONPath: "properties.clusterVersion", Default: true},
	{Name: "SnapshotSize", JSONPath: "properties.snapshotSize"},
	{Name: "RequiredSizeForRestore", JSONPath: "properties.requiredSizeForRestore"},
	{Name: "EarliestRecoveryTargetTime", JSONPath: "properties.earliestRecoveryTargetTime"},
	{Name: "LatestRecoveryTargetTime", JSONPath: "properties.latestRecoveryTargetTime"},
}

func SnapshotCmd() *core.Command {
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "snapshot",
			Aliases: []string{"snap", "snaps", "snapshots"},
			Short:   "In-Memory DB Snapshot Operations",
			Long: `The sub-commands of ` + "`ionosctl dbaas in-memory-db-v2 snapshot`" + ` allow you to view In-Memory DB Cluster Snapshots.

A snapshot is not a single point in time — it represents a recovery WINDOW between its earliestRecoveryTargetTime and latestRecoveryTargetTime. Backups combine periodic dumps with a continuous change log, so a cluster can be restored to any moment within that window. Use ` + "`--recovery-time`" + ` on ` + "`cluster restore`" + ` (or ` + "`cluster create --snapshot-id`" + `) to zoom into a specific point in the window; omit it to use the latest point.`,
			TraverseChildren: true,
		},
	}

	snapshotCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(snapshotCols))
	_ = snapshotCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(snapshotCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	snapshotCmd.AddCommand(SnapshotListCmd())
	snapshotCmd.AddCommand(SnapshotGetCmd())
	snapshotCmd.AddCommand(location.LocationCmd())

	return snapshotCmd
}
