package snapshot

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "snapshot",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List the snapshots of your Mongo Cluster",
		Example:   "ionosctl dbaas mongo cluster snapshot ls --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Printer.Verbose("Getting snapshots of Cluster %s", clusterId)
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
			snapshots, _, err := c.DbaasMongoServices.Clusters().SnapshotsList(clusterId, limitPtr, offsetPtr)
			if err != nil {
				return err
			}
			return c.Printer.Print(getSnapshotPrint(c, snapshots.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
