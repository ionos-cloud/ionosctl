package database

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Owner", JSONPath: "properties.owner", Default: true},
	{Name: "ClusterId", JSONPath: "ClusterId"},
}

func DatabaseCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "database",
			Aliases: []string{"databases"},
			Short:   "Manage logical databases inside a PostgreSQL cluster",
			Long: `Manage the logical PostgreSQL databases hosted inside a cluster.

A database is a named logical database within a cluster, owned by exactly one user (role). The owner is granted full privileges on that database, so the owner you assign must be an existing user in the same cluster (see 'dbaas postgres user'). A single cluster can host many databases.

Databases can only be created once the cluster is AVAILABLE.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DeleteCmd())
	return cmd
}
