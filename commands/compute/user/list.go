package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

func UserListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "user",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List all Users on the contract",
		LongDesc:   "List every User on your contract. The output shows each User's admin status, whether two-factor auth is forced/active, and whether the account is active (not blacklisted).\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.UsersFiltersUsage(),
		Example:    "ionosctl compute user list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunUserList,
		InitClient: true,
	})

	cmd.Command.Flags().StringSlice("cols", nil, table.ColsMessage(allUserCols))
	_ = cmd.Command.RegisterFlagCompletionFunc("cols", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allUserCols), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
