package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterListCmd() *core.Command {
	ctx := context.TODO()
	listEnv := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"ls", "l"},
		ShortDesc:  "List MariaDB Clusters",
		LongDesc:   "Use this command to retrieve a list of MariaDB Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:    "ionosctl dbaas mariadb-v2 cluster list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})
	listEnv.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the MariaDB Clusters that contain the specified name in the Name field. The value is case insensitive")
	listEnv.AddStringFlag(constants.FlagState, "", "", "Response filter by cluster state: PROVISIONING, AVAILABLE, UPDATING, DESTROYING, FAILED")
	_ = listEnv.Command.RegisterFlagCompletionFunc(constants.FlagState, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"PROVISIONING", "AVAILABLE", "UPDATING", "DESTROYING", "FAILED"}, cobra.ShellCompDirectiveNoFileComp
	})
	listEnv.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	listEnv.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return listEnv
}

func RunClusterList(c *core.CommandConfig) error {
	return c.ListAllLocations(clusterCols, func(cfg *shared.Configuration) (any, error) {
		apiClient := mariadb.NewAPIClient(cfg)
		req := apiClient.ClustersApi.ClustersGet(context.Background())

		if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
			req = req.FilterName(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
			req = req.FilterState(mariadb.MariadbClusterStates(viper.GetString(fn)))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
			req = req.Limit(viper.GetInt32(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
			req = req.Offset(viper.GetInt32(fn))
		}

		clusters, _, err := req.Execute()
		return clusters, err
	})
}
