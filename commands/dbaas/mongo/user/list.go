package user

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func listAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Users from all clusters...")
	ls, err := c.DbaasMongoServices.Users().ListAll()
	if err != nil {
		return err
	}
	return c.Printer.Print(getUserPrint(c, &ls))
}

func UserListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "Retrieves a list of MongoDB users.",
		Example:   "ionosctl dbaas mongo user list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return listAll(c)
			}
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Printer.Verbose("Getting Users from all cluster %s", clusterId)

			ls, _, err := c.DbaasMongoServices.Users().List(clusterId)
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(c, ls.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all users, across all clusters")
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
