package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Username", JSONPath: "properties.username", Default: true},
	{Name: "System", JSONPath: "properties.system", Default: true},
	{Name: "ClusterId", JSONPath: "ClusterId"},
}

func UserCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "user",
			Aliases: []string{"usr", "u", "users"},
			Short:   "Manage PostgreSQL users (roles)",
			Long: `Manage the PostgreSQL users (login roles) of a cluster.

A user is a PostgreSQL role that can log in with a password. The initial user is created together with the cluster ('cluster create --db-username/--db-password'); this command manages additional users. A user can be made the owner of a database (see 'dbaas postgres database create --owner').

System users (created and managed by the service, e.g. postgres) are shown with System=true; they cannot be updated or deleted here, and 'user list' hides them unless you pass --system.

Users can only be created once the cluster has reached the AVAILABLE state.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(UpdateCmd())
	return cmd
}
