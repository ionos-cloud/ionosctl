package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/snapshot/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
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
			Use:              "snapshot",
			Aliases:          []string{"snap", "snaps", "snapshots"},
			Short:            "In-Memory DB Snapshot Operations",
			Long:             "The sub-commands of `ionosctl dbaas in-memory-db-v2 snapshot` allow you to manage In-Memory DB Cluster Snapshots.",
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

// Snapshots returns all snapshots matching the given filters
func Snapshots(fs ...Filter) (inmemorydb.SnapshotReadList, error) {
	req := client.Must().InMemoryDBClientV2.SnapshotsApi.SnapshotsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return inmemorydb.SnapshotReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return inmemorydb.SnapshotReadList{}, err
	}
	return ls, nil
}

type Filter func(request inmemorydb.ApiSnapshotsGetRequest) (inmemorydb.ApiSnapshotsGetRequest, error)
