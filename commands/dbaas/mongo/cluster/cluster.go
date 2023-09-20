package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
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
	clusterCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

var (
	allJSONPaths = map[string]string{
		"ClusterId":    "id",
		"Name":         "properties.displayName",
		"Edition":      "properties.edition",
		"Type":         "properties.type",
		"URL":          "properties.connectionString",
		"Instances":    "properties.instances",
		"Shards":       "properties.shards",
		"Health":       "metadata.health",
		"State":        "metadata.state",
		"MongoVersion": "properties.mongoDBVersion",
		"Location":     "properties.location",
		"TemplateId":   "properties.templateID",
		"Cores":        "properties.cores",
		"StorageType":  "properties.storageType",
	}

	allCols = []string{"ClusterId", "Name", "Edition", "Type", "URL", "Instances", "Shards", "Health", "State",
		"MongoVersion", "MaintenanceWindow", "Location", "DatacenterId", "LanId", "Cidr", "TemplateId", "Cores", "RAM",
		"StorageSize", "StorageType"}

	defaultCols = allCols[0:9]
)

func convertClusterToTable(cluster ionoscloud.ClusterResponse) ([]map[string]interface{}, error) {
	properties, ok := cluster.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster properties")
	}

	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
	if !ok || maintenanceWindow == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window")
	}

	day, ok := maintenanceWindow.GetDayOfTheWeekOk()
	if !ok || day == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window day")
	}

	tyme, ok := maintenanceWindow.GetTimeOk()
	if !ok || tyme == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window time")
	}

	storage, ok := properties.GetStorageSizeOk()
	if !ok || storage == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster storage size")
	}

	ram, ok := properties.GetRamOk()
	if !ok || ram == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Cluster RAM")
	}

	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, cluster)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)
	temp[0]["RAM"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*ram), convbytes.MB, convbytes.GB))
	temp[0]["StorageSize"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*storage), convbytes.MB, convbytes.GB))

	connections, ok := properties.GetConnectionsOk()
	if ok && connections != nil {
		for _, con := range *connections {
			dcId, ok := con.GetDatacenterIdOk()
			if !ok || dcId == nil {
				return nil, fmt.Errorf("could not retrieve Mongo Cluster datacenter ID")
			}

			lanId, ok := con.GetLanIdOk()
			if !ok || lanId == nil {
				return nil, fmt.Errorf("could not retrieve Mongo Cluster lan ID")
			}

			cidr, ok := con.GetCidrListOk()
			if !ok || cidr == nil {
				return nil, fmt.Errorf("could not retrieve Mongo Cluster CIDRs")
			}

			temp[0]["DatacenterId"] = *dcId
			temp[0]["LanId"] = *lanId
			temp[0]["Cidr"] = *cidr
		}
	}
	return temp, nil
}

func convertClustersToTable(clusters ionoscloud.ClusterList) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := convertClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func Clusters(fs ...Filter) (ionoscloud.ClusterList, error) {
	req := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		req = f(req)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return ionoscloud.ClusterList{}, fmt.Errorf("failed getting clusters: %w", err)
	}
	return clusters, err
}

type Filter func(ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
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
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
