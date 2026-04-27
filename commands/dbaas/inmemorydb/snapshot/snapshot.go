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
			Use:              "snapshot",
			Aliases:          []string{"snaps", "snap", "backup", "backups", "snapshots"},
			Short:            "The sub-commands of 'ionosctl dbaas inmemorydb snapshots' allow you to manage In-Memory DB Replica Set Snapshots.",
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(List())
	cmd.AddCommand(restore.Root())

	return cmd
}
