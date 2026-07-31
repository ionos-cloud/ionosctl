package restore

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.displayName", Default: true},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "ReplicasetId", JSONPath: "properties.replicasetId", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "RestoreTime", JSONPath: "metadata.restoreTime", Default: true},
	{Name: "RestoredSnapshotId", JSONPath: "metadata.restoredSnapshotId", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "restore",
			Aliases: []string{"restores", "backup", "backups"},
			Short:   "Manage In-Memory DB Restores (restoring a replica set from a snapshot)",
			Long: `The sub-commands of ` + "`ionosctl dbaas in-memory-db snapshot restore`" + ` manage Restores.

A restore takes an existing point-in-time snapshot and applies it to a target replica set (--replicaset-id), rolling that replica set's data back to the snapshot's state. Each restore is itself a tracked resource with its own state and RestoreTime, so you can list past restore operations for a snapshot.

The snapshot and the target replica set must be in the same location/datacenter (snapshots are not portable across datacenters). To instead create a brand-new replica set from a snapshot, use 'replicaset create --snapshot-id'.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Create())
	cmd.AddCommand(List())

	return cmd
}
