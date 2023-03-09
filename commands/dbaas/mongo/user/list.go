package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
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
			var limitPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
				limit := viper.GetInt32(f)
				limitPtr = &limit
			}
			var offsetPtr *int32 = nil
			if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
				offset := viper.GetInt32(f)
				offsetPtr = &offset
			}
			ls, _, err := c.DbaasMongoServices.Users().List(clusterId, limitPtr, offsetPtr)
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
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all users, across all clusters")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	cmd.Command.SilenceUsage = true

	cmd.Command.SilenceUsage = true

	return cmd
}
