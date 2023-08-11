package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func confirmStringForCluster(c sdkgo.ClusterResponse) string {
	askString := ""
	if p := c.Properties; p != nil {
		if edition := p.Edition; edition != nil {
			askString = fmt.Sprintf("%s %s", askString, *edition)
		}
		if ctype := p.Type; ctype != nil {
			askString = fmt.Sprintf("%s %s", askString, *ctype)
		}
		if c.Id != nil {
			askString = fmt.Sprintf("%s cluster %s", askString, *c.Id)
		}
		if n := p.DisplayName; n != nil {
			askString = fmt.Sprintf("%s (%s)", askString, *n)
		}
		if v := p.MongoDBVersion; v != nil {
			askString = fmt.Sprintf("%s version v%s", askString, *v)
		}
		if l := p.Location; l != nil {
			askString = fmt.Sprintf("%s located in %s", askString, *l)
		}
	}
	return fmt.Sprintf("delete%s and its snapshots", askString)
}

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
			chosenCluster, _, err := client.Must().MongoClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				wrapped := fmt.Errorf("failed trying to find cluster by id: %w", err)
				keepGoing := confirm.Ask(fmt.Sprintf("%s, try deleting %s anyways", wrapped.Error(), clusterId))
				if !keepGoing {
					return wrapped
				}
			}

			ok := confirm.Ask(confirmStringForCluster(chosenCluster), viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
			if !ok {
				return fmt.Errorf("user denied confirmation")
			}
			c.Printer.Verbose("Deleting cluster: %s", clusterId)

			_, _, err = client.Must().MongoClient.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
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

func deleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Deleting All Clusters!")
	xs, _, err := client.Must().MongoClient.ClustersApi.ClustersGet(c.Context).Execute()
	if err != nil {
		return err
	}

	return functional.ApplyAndAggregateErrors(*xs.GetItems(), func(x sdkgo.ClusterResponse) error {
		yes := confirm.Ask(confirmStringForCluster(x), viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
		if yes {
			_, _, delErr := client.Must().MongoClient.ClustersApi.ClustersDelete(c.Context, *x.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting cluster %s: %w", *x.Properties.DisplayName, delErr)
			}
		}
		return nil
	})
}
