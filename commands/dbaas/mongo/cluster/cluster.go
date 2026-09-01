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
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "Mongo Cluster Operations",
			Long: `The sub-commands of ` + "`ionosctl dbaas mongo cluster`" + ` manage MongoDB clusters: managed, IONOS-hosted MongoDB deployments.

A cluster is defined by three things that together determine its shape and price:
  - Edition — the tier that bounds the allowed sizing and topology:
      * playground  - a single-instance sandbox (1 instance, fixed 1 core / 2 GB RAM / 50 GB storage via the Playground template). No snapshots, not for production.
      * business    - a replicaset for typical production workloads. Sized via a template (XS...4XL) that bundles cores/RAM/storage. Daily snapshots retained 7 days.
      * enterprise  - a replicaset OR a sharded-cluster with explicitly chosen cores/RAM/storage. Adds point-in-time restore (recover to any moment in the last ~7 days).
  - Type — 'replicaset' (one primary + n-1 secondaries, all holding the same data) or 'sharded-cluster' (data partitioned across shards, each shard itself a replicaset). Only enterprise clusters may be sharded.
  - Sizing — either a template (playground/business) or explicit --cores/--ram/--storage-size/--storage-type (enterprise).

Instances are the MongoDB nodes of a replicaset (odd counts 1/3/5/7 so a majority can elect a primary). Shards partition data horizontally in a sharded-cluster (2-32). Clients reach the cluster over a private LAN in a datacenter (--datacenter-id, --lan-id, --cidr); the connection string is shown as the URL column. Maintenance happens in a weekly 4-hour window (--maintenance-day / --maintenance-time). Backups live in IONOS S3 Object Storage; see 'snapshot' and 'cluster restore'.`,
			TraverseChildren: true,
		},
	}

	clusterCmd.AddColsFlag(allCols)

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
