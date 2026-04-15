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
			Use:              "restore",
			Aliases:          []string{"restores", "backup", "backups"},
			Short:            "The sub-commands of 'ionosctl dbaas inmemorydb restore' allow you to manage In-Memory DB Replica Set Restores.",
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Create())
	cmd.AddCommand(List())

	return cmd
}
