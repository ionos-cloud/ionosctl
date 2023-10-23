package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsonpaths"
	"github.com/ionos-cloud/sdk-go-dataplatform"
)

func ConvertDataplatformClusterToTable(cluster ionoscloud.ClusterResponseData) ([]map[string]interface{}, error) {
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

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DataplatformCluster, cluster)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)

	return temp, nil
}

func ConvertDataplatformClustersToTable(clusters ionoscloud.ClusterListResponseData) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Dataplatform Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertDataplatformClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertDataplatformNodePoolToTable(np ionoscloud.NodePoolResponseData) ([]map[string]interface{}, error) {
	properties, ok := np.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool properties")
	}

	ramRaw, ok := properties.GetRamSizeOk()
	if !ok || ramRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool RAM size")
	}

	gb := convbytes.Convert(int64(*ramRaw), convbytes.MB, convbytes.GB)
	ram := fmt.Sprintf("%v GB", gb)

	storageSizeRaw, ok := properties.GetStorageSizeOk()
	if !ok || storageSizeRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Storage size")
	}

	storageTypeRaw, ok := properties.GetStorageTypeOk()
	if !ok || storageTypeRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Storage type")
	}

	storageGb := convbytes.Convert(int64(*storageSizeRaw), convbytes.MB, convbytes.GB)
	storage := fmt.Sprintf("%v %v GB", *storageTypeRaw, storageGb)

	maintenanceWindowRaw, ok := properties.GetMaintenanceWindowOk()
	if !ok || maintenanceWindowRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window")
	}

	day, ok := maintenanceWindowRaw.GetDayOfTheWeekOk()
	if !ok || day == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window Day")
	}

	time, ok := maintenanceWindowRaw.GetTimeOk()
	if !ok || time == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window Time")
	}

	maintenanceWindow := fmt.Sprintf("%v %v", *day, *time)

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DataplatformNodepool, np)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["Ram"] = ram
	temp[0]["Storage"] = storage
	temp[0]["MaintenanceWindow"] = maintenanceWindow

	return temp, nil
}

func ConvertDataplatformNodePoolsToTable(nps ionoscloud.NodePoolListResponseData) ([]map[string]interface{}, error) {
	var npsConverted = make([]map[string]interface{}, 0)

	items, ok := nps.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Items")
	}

	for _, item := range *items {
		temp, err := ConvertDataplatformNodePoolToTable(item)
		if err != nil {
			return nil, err
		}

		npsConverted = append(npsConverted, temp...)
	}

	return npsConverted, nil
}
