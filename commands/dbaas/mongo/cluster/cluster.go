package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	enumEditions = []string{"playground", "business", "enterprise"} // Remove whenever the SDK adds this as an actual type with enum vals
	enumTypes    = []string{"replicaset", "sharded-cluster"}        // Remove whenever the SDK adds this as an actual type with enum vals
)

const (
	flagBackupLocation     = "backup-location"
	flagBiconnector        = "biconnector"
	flagBiconnectorEnabled = "biconnector-enabled"
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Mongo Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas mongo cluster` allow you to manage the Mongo Clusters under your account.",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

var (
	allCols = []string{"ClusterId", "Name", "Edition", "Type", "URL", "Instances", "Shards", "Health", "State",
		"MongoVersion", "MaintenanceWindow", "Location", "DatacenterId", "LanId", "Cidr", "TemplateId", "Cores", "RAM",
		"StorageSize", "StorageType"}

	defaultCols = allCols[0:9]
)

func Clusters(fs ...Filter) (mongo.ClusterList, error) {
	req := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		req = f(req)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return mongo.ClusterList{}, fmt.Errorf("failed getting clusters: %w", err)
	}
	return clusters, err
}

type Filter func(mongo.ApiClustersGetRequest) mongo.ApiClustersGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req mongo.ApiClustersGetRequest) mongo.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
			req = req.Limit(viper.GetInt32(f))
		}
		if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
			req = req.Offset(viper.GetInt32(f))
		}
		return req
	}
}

func FilterNameFlags(c *core.CommandConfig) Filter {
	return func(req mongo.ApiClustersGetRequest) mongo.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
