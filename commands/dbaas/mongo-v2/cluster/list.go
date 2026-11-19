package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterListCmd() *core.Command {
	ctx := context.TODO()
	listEnv := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mongo-v2",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"ls", "l"},
		ShortDesc:  "List MongoDB Clusters",
		LongDesc:   "Use this command to retrieve a list of MongoDB Clusters provisioned under your account. You can filter the result by exact Cluster Name using `--name`, by state using `--state` and by edition using `--edition`.",
		Example:    "ionosctl dbaas mongo-v2 cluster list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})
	listEnv.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the MongoDB Clusters with exactly this name")
	listEnv.AddStringFlag(constants.FlagState, "", "", "Response filter by cluster state: PROVISIONING, AVAILABLE, UPDATING, DESTROYING, FAILED")
	_ = listEnv.Command.RegisterFlagCompletionFunc(constants.FlagState, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"PROVISIONING", "AVAILABLE", "UPDATING", "DESTROYING", "FAILED"}, cobra.ShellCompDirectiveNoFileComp
	})
	listEnv.AddSetFlag(constants.FlagEdition, "", "", editions, "Response filter to list only the MongoDB Clusters of this edition")
	listEnv.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	listEnv.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return listEnv
}

func RunClusterList(c *core.CommandConfig) error {
	return c.ListAllLocations(clusterCols, func(cfg *shared.Configuration) (any, error) {
		apiClient := mongo.NewAPIClient(cfg)
		req := apiClient.ClustersApi.ClustersGet(context.Background())

		if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
			req = req.FilterName(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
			req = req.FilterState(mongo.MongoClusterStates(viper.GetString(fn)))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fn) {
			req = req.FilterEdition(mongo.MongoEdition(viper.GetString(fn)))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
			req = req.Limit(viper.GetInt32(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
			req = req.Offset(viper.GetInt32(fn))
		}

		clusters, _, err := req.Execute()
		if err != nil {
			return clusters, err
		}
		if fn := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fn) {
			edition := mongo.MongoEdition(viper.GetString(fn))
			clusters.Items = functional.Filter(clusters.Items, func(cl mongo.ClusterRead) bool {
				return clusterEdition(cl.Properties) == edition
			})
		}
		return clusters, nil
	})
}

func clusterEdition(cluster mongo.Cluster) mongo.MongoEdition {
	if cluster.BusinessCluster != nil {
		return mongo.MONGOEDITION_BUSINESS
	}
	return mongo.MONGOEDITION_ENTERPRISE
}
