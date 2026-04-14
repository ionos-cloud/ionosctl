package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
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
			Use:              "user",
			Aliases:          []string{"usr", "u", "users"},
			Short:            "DBaaS Postgresql User Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres user` allow you to perform operations on DBaaS PostgreSQL users.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(UpdateCmd())
	return cmd
}
