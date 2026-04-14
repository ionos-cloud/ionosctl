package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
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

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

var allCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.displayName", Default: true},
	{Name: "Edition", JSONPath: "properties.edition", Default: true},
	{Name: "Type", JSONPath: "properties.type", Default: true},
	{Name: "URL", JSONPath: "properties.connectionString", Default: true},
	{Name: "Instances", JSONPath: "properties.instances", Default: true},
	{Name: "Shards", JSONPath: "properties.shards", Default: true},
	{Name: "Health", JSONPath: "metadata.health", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "MongoVersion", JSONPath: "properties.mongoDBVersion"},
	{Name: "MaintenanceWindow", Format: func(item map[string]any) any {
		day, _ := table.Navigate(item, "properties.maintenanceWindow.dayOfTheWeek").(string)
		t, _ := table.Navigate(item, "properties.maintenanceWindow.time").(string)
		if day == "" && t == "" {
			return nil
		}
		return fmt.Sprintf("%s %s", day, t)
	}},
	{Name: "Location", JSONPath: "properties.location"},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId"},
	{Name: "Cidr", JSONPath: "properties.connections.0.cidrList"},
	{Name: "TemplateId", JSONPath: "properties.templateID"},
	{Name: "Cores", JSONPath: "properties.cores"},
	{Name: "RAM", Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.ram")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d GB", int(f/1024))
	}},
	{Name: "StorageSize", Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.storageSize")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d GB", int(f/1024))
	}},
	{Name: "StorageType", JSONPath: "properties.storageType"},
}

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

func FilterNameFlags(c *core.CommandConfig) Filter {
	return func(req mongo.ApiClustersGetRequest) mongo.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
