package nodepool

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"
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
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId, constants.FlagNodepoolId}, []string{constants.ArgAll, constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c, clusterId)
			}

			nodepoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
			np, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsFindById(context.Background(), clusterId, nodepoolId).Execute()
			if err != nil {
				return fmt.Errorf("couldn't find nodepool: %w", err)
			}

			ok := confirm.Ask(fmt.Sprintf("delete nodepool %s (%s)", nodepoolId, *np.Properties.Name), viper.GetBool(constants.ArgForce))
			if !ok {
				return fmt.Errorf("canceled deletion: invalid input")
			}

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting nodepool: %s", nodepoolId))
			_, _, err = client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(c.Context, clusterId, nodepoolId).Execute()
			return err
		},
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, constants.FlagIdShort, "", "The unique ID of the nodepool", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformNodepoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all clusters. If cluster ID is provided, delete all nodepools in given cluster")

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig, clusterId string) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all nodepools!"))
	if clusterId != "" {
		// Both --all and --cluster-id provided, so delete only the nodepools of the given cluster.
		return deleteNodePools(c, clusterId)
	}
	// Only --all is provided, so delete all nodepools

	ls, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(c.Context).Execute()
	if err != nil {
		return err
	}

	// accumulate the error. If it's not nil break out of the fold
	return functional.ApplyAndAggregateErrors(*ls.GetItems(), func(x ionoscloud.ClusterResponseData) error {
		return deleteNodePools(c, *x.Id)
	})
}

func deleteNodePools(c *core.CommandConfig, clusterId string) error {
	xs, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}
	return functional.ApplyAndAggregateErrors(*xs.GetItems(), func(x ionoscloud.NodePoolResponseData) error {
		ok := confirm.Ask(fmt.Sprintf("delete nodepool %s (%s)", *x.Id, *x.Properties.Name), viper.GetBool(constants.ArgForce))
		if !ok {
			return fmt.Errorf("canceled deletion for %s (%s): invalid input", *x.Id, *x.Properties.Name)
		}
		_, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsDelete(context.Background(), clusterId, *x.Id).Execute()
		return err
	})
}
