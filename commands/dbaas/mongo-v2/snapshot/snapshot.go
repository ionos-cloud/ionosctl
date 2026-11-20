package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/snapshot/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var snapshotCols = []table.Column{
	{Name: "SnapshotId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "ClusterName", JSONPath: "properties.clusterName", Default: true},
	{Name: "Version", JSONPath: "properties.mongoVersion", Default: true},
	{Name: "SnapshotLocation", JSONPath: "properties.location", Default: true},
	{Name: "EarliestRecoveryTargetTime", JSONPath: "properties.earliestRecoveryTargetTime", Default: true},
	{Name: "LatestRecoveryTargetTime", JSONPath: "properties.latestRecoveryTargetTime"},
	{Name: "TotalSize", JSONPath: "properties.totalSize"},
	{Name: "RequiredStorageSize", JSONPath: "properties.requiredStorageSize"},
}

func SnapshotCmd() *core.Command {
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "snapshot",
			Aliases: []string{"snap", "snaps", "snapshots"},
			Short:   "MongoDB Snapshot Operations",
			Long: `The sub-commands of ` + "`ionosctl dbaas mongo-v2 snapshot`" + ` allow you to view MongoDB Cluster Snapshots.

A snapshot is not a single point in time — it keeps the cluster's oplog, so it represents a recovery WINDOW from its earliestRecoveryTargetTime to its latestRecoveryTargetTime (empty while the window extends to the present). A cluster can be restored to any moment within that window. Use ` + "`--recovery-time`" + ` on ` + "`cluster restore`" + ` (or ` + "`cluster create --snapshot-id`" + `) to pick a specific point in the window. TotalSize and RequiredStorageSize are in MiB.`,
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
