package cluster

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
)

var (
	allJSONPaths = map[string]string{
		"Id":           "id",
		"Name":         "properties.name",
		"Version":      "properties.dataPlatformVersion",
		"DatacenterId": "properties.datacenterId",
		"State":        "metadata.state",
	}

	allCols = []string{"Id", "Name", "Version", "MaintenanceWindow", "DatacenterId", "State"}
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Dataplatform Cluster Operations",
			Long:             "This command allows you to interact with the already created clusters or creates new clusters in your virtual data center",
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
	clusterCmd.AddCommand(ClustersKubeConfigCmd())

	return clusterCmd
}

func convertClusterToTable(cluster ionoscloud.ClusterResponseData) ([]map[string]interface{}, error) {
	properties, ok := cluster.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Cluster properties")
	}

	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
	if !ok || maintenanceWindow == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Cluster maintenance window")
	}

	day, ok := maintenanceWindow.GetDayOfTheWeekOk()
	if !ok || day == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Cluster maintenance window day")
	}

	tyme, ok := maintenanceWindow.GetTimeOk()
	if !ok || tyme == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Cluster maintenance window time")
	}

	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, cluster)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)

	return temp, nil
}

func convertClustersToTable(clusters ionoscloud.ClusterListResponseData) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Clusters items")
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
