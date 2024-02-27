package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

func GetCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "Get user",
			LongDesc:  "Get the specified user in the given database cluster",
			Example:   "ionosctl dbaas-postgres user get --cluster-id <cluster-id> --user <user>",
			PreCmdRun: core.NoPreRun,
			CmdRun:    RunCmd,
		},
	)
	c.Command.Flags().StringSlice(constants.ArgCols, []string{}, tabheaders.ColsMessage(allCols))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveDefault
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveDefault
		},
	)

	c.AddStringFlag(constants.ArgUser, "", "", "The name of the user to retrieve")
	_ = c.Command.RegisterFlagCompletionFunc(constants.ArgUser, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UserNames(c), cobra.ShellCompDirectiveDefault
	})

	c.AddBoolFlag("system", "", false, "List system users along with regular users")

	return c
}
