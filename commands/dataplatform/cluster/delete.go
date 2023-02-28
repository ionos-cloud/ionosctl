package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func deleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Deleting All Clusters!")
	if !viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)) {
		err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all clusters")
		if err != nil {
			return err
		}
	}
	_, err := c.DbaasMongoServices.Clusters().DeleteAll(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	return err
}

func ClusterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a Dataplatform Cluster by ID",
		Example:   "ionosctl dataplatform cluster delete --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, []string{constants.ArgAll}, []string{constants.FlagClusterId})
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
			client, err := config.GetClient()
			_, _, err = client.DataplatformClient.DataPlatformClusterApi.DeleteCluster(c.Context, clusterId).Execute()
			if err != nil {
				return err
			}
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all clusters")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no verification")

	cmd.Command.SilenceUsage = true

	return cmd
}
