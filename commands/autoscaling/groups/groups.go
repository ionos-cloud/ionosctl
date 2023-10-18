package groups

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
)

func GroupsCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "groups",
			Aliases:          []string{"g", "group"},
			Short:            "Autoscaling Groups Operations",
			Long:             "The sub-commands of `ionosctl autoscaling groups` allow you to manage the Autoscaling Groups under your account.",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	clusterCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	return clusterCmd
}

var (
	allJSONPaths = map[string]string{
		"GroupId":      "id",
		"Name":         "properties.name",
		"MinReplicas":  "properties.minReplicaCount",
		"MaxReplicas":  "properties.maxReplicaCount",
		"DatacenterId": "properties.datacenter.id",
		"State":        "metadata.state",
	}

	allCols     = []string{"GroupId", "Name", "MinReplicas", "MaxReplicas", "DatacenterId", "State"}
	defaultCols = allCols[0:9]
)

// func toTable(cluster ionoscloud.ClusterResponse) ([]map[string]interface{}, error) {
// 	properties, ok := cluster.GetPropertiesOk()
// 	if !ok || properties == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster properties")
// 	}
//
// 	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
// 	if !ok || maintenanceWindow == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window")
// 	}
//
// 	day, ok := maintenanceWindow.GetDayOfTheWeekOk()
// 	if !ok || day == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window day")
// 	}
//
// 	tyme, ok := maintenanceWindow.GetTimeOk()
// 	if !ok || tyme == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster maintenance window time")
// 	}
//
// 	storage, ok := properties.GetStorageSizeOk()
// 	if !ok || storage == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster storage size")
// 	}
//
// 	ram, ok := properties.GetRamOk()
// 	if !ok || ram == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Cluster RAM")
// 	}
//
// 	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, cluster)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
// 	}
//
// 	temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)
// 	temp[0]["RAM"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*ram), convbytes.MB, convbytes.GB))
// 	temp[0]["StorageSize"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*storage), convbytes.MB, convbytes.GB))
//
// 	connections, ok := properties.GetConnectionsOk()
// 	if ok && connections != nil {
// 		for _, con := range *connections {
// 			dcId, ok := con.GetDatacenterIdOk()
// 			if !ok || dcId == nil {
// 				return nil, fmt.Errorf("could not retrieve Mongo Cluster datacenter ID")
// 			}
//
// 			lanId, ok := con.GetLanIdOk()
// 			if !ok || lanId == nil {
// 				return nil, fmt.Errorf("could not retrieve Mongo Cluster lan ID")
// 			}
//
// 			cidr, ok := con.GetCidrListOk()
// 			if !ok || cidr == nil {
// 				return nil, fmt.Errorf("could not retrieve Mongo Cluster CIDRs")
// 			}
//
// 			temp[0]["DatacenterId"] = *dcId
// 			temp[0]["LanId"] = *lanId
// 			temp[0]["Cidr"] = *cidr
// 		}
// 	}
// 	return temp, nil
// }

// func lsToTable(ls ionoscloud.GroupsList) ([]map[string]interface{}, error) {
// 	items, ok := clusters.GetItemsOk()
// 	if !ok || items == nil {
// 		return nil, fmt.Errorf("could not retrieve Mongo Clusters items")
// 	}
//
// 	var clustersConverted []map[string]interface{}
// 	for _, item := range *items {
// 		temp, err := convertClusterToTable(item)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		clustersConverted = append(clustersConverted, temp...)
// 	}
//
// 	return clustersConverted, nil
// }
