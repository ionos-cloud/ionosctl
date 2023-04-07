package nodepool

import (
	"context"
	"fmt"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a Dataplatform Cluster by ID",
		Example:   "ionosctl dataplatform cluster delete --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId}, []string{constants.ArgAll, constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c, clusterId)
			}

			nodepoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
			err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete nodepool %s", nodepoolId))
			if err != nil {
				return err
			}
			c.Printer.Verbose("Deleting nodepool: %s", nodepoolId)
			client, err := client2.Get()
			_, _, err = client.DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(c.Context, clusterId, nodepoolId).Execute()
			if err != nil {
				return err
			}
			return err
		},
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, constants.FlagIdShort, "", "The unique ID of the nodepool", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all clusters. If cluster ID is provided, delete all nodepools in given cluster")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no verification")

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig, clusterId string) error {
	c.Printer.Verbose("Deleting all nodepools!")
	if !viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)) {
		err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all nodepools")
		if err != nil {
			return err
		}
	}
	if clusterId != "" {
		// Both --all and --cluster-id provided, so delete only the nodepools of the given cluster.
		return deleteNodePools(clusterId)
	}
	// Only --all is provided, so delete all nodepools

	client, err := client2.Get()
	if err != nil {
		return err
	}
	ls, _, err := client.DataplatformClient.DataPlatformClusterApi.ClustersGet(c.Context).Execute()
	if err != nil {
		return err
	}

	// accumulate the error. If it's not nil break out of the fold
	return shared.ApplyOrFail(*ls.GetItems(), func(x ionoscloud.ClusterResponseData) error {
		return deleteNodePools(*x.Id)
	})
}

func deleteNodePools(clusterId string) error {
	client, err := client2.Get()
	if err != nil {
		return err
	}
	xs, _, err := client.DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}
	return shared.ApplyOrFail(*xs.GetItems(), func(x ionoscloud.NodePoolResponseData) error {
		_, _, err := client.DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(context.Background(), clusterId, *x.Id).Execute()
		return err
	})
}
