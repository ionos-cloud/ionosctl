package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/snapshot/restore"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "ReplicasetId", JSONPath: "metadata.replicasetId", Default: true},
	{Name: "DatacenterId", JSONPath: "metadata.datacenterId", Default: true},
	{Name: "Time", JSONPath: "metadata.snapshotTime", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "snapshot",
			Aliases: []string{"snaps", "snap", "backup", "backups", "snapshots"},
			Short:   "Manage In-Memory DB Snapshots (point-in-time backups)",
			Long: `The sub-commands of ` + "`ionosctl dbaas in-memory-db snapshot`" + ` let you view In-Memory DB Snapshots and restore from them.

A snapshot is an automatic, read-only, point-in-time dump of a replica set (see the Time column for when it was taken). Snapshots are stored per datacenter and are only usable within the location/datacenter of the replica set they belong to. Snapshots cannot be created or deleted manually here - the platform takes them for you.

Use them in two ways:
  - 'snapshot restore create': restore a snapshot onto an existing replica set (--replicaset-id), rolling that replica set back to the snapshot's state.
  - 'replicaset create --snapshot-id': create a brand-new replica set restored from a snapshot.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(restore.Root())

	return cmd
}
