package cluster

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all mongo clusters")
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no verification")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	client, err := config.GetClient()
	if err != nil {
		return err
	}
	c.Printer.Verbose("Deleting All Clusters!")
	xs, _, err := client.MongoClient.ClustersApi.ClustersGet(c.Context).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyOrFail(*xs.GetItems(), func(x sdkgo.ClusterResponse) error {
		yes := confirm.Ask(fmt.Sprintf("delete cluster %s (%s)", *x.Id, *x.Properties.DisplayName), viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
		if yes {
			_, _, delErr := client.MongoClient.ClustersApi.ClustersDelete(c.Context, *x.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting one of the clusters: %w", delErr)
			}
		}
		return nil
	})
}
