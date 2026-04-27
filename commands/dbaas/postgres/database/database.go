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
			Use:              "database",
			Aliases:          []string{"databases"},
			Short:            "DBaaS Postgresql Database Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres database` allow you to perform operations on DBaaS PostgreSQL databases.",
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
