package cluster

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/viper"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
)

func deleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Deleting All Clusters!")
	if !viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)) {
		// TODO: This is a pretty nasty snippet to duplicate everywhere
		err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all clusters")
		if err != nil {
			return err
		}
	}
	_, err := c.DbaasMongoServices.Clusters().DeleteAll(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	return err
}

func ClusterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a Mongo Cluster by ID",
		Example:   "ionosctl dbaas mongo cluster delete --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete cluster %s", clusterId))
			if err != nil {
				return err
			}
			c.Printer.Verbose("Deleting cluster: %s", clusterId)
			_, err = c.DbaasMongoServices.Clusters().Delete(clusterId)
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all mongo clusters")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no verification")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
