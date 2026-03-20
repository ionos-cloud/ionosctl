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
	{Name: "ClusterId"},
}

func UserCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"usr", "u", "users"},
			Short:            "DBaaS Postgresql User Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres user` allow you to perform operations on DBaaS PostgreSQL users.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(UpdateCmd())
	return cmd
}
