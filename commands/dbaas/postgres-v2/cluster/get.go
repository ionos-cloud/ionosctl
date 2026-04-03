package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Cluster",
		Example:    "ionosctl dbaas postgres-v2 cluster get --cluster-id <cluster-id>",
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})

	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			clusters, err := Clusters()
			if err != nil {
				return []string{}
			}
			return functional.Map(clusters.Items, func(c psqlv2.ClusterRead) string {
				return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
					c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)

			})
		}, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state [seconds]")
	get.AddStringSliceFlag(constants.ArgCols, "", table.DefaultCols(clusterCols), table.ColsMessage(clusterCols))
	_ = get.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(clusterCols), cobra.ShellCompDirectiveNoFileComp
	})
	return get
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterGet(c *core.CommandConfig) error {
	c.Verbose(constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	c.Verbose("Getting Cluster...")

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err := waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
	}

	cluster, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(
		context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).Execute()
	if err != nil {
		return fmt.Errorf("could not get cluster: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}
