package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
)

func ConvertK8sClusterToTable(cluster ionoscloud.KubernetesCluster) ([]map[string]interface{}, error) {
	properties, ok := cluster.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve K8s Cluster properties")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.K8sCluster, cluster)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
	if ok && maintenanceWindow != nil {
		day, ok := maintenanceWindow.GetDayOfTheWeekOk()
		if !ok || day == nil {
			return nil, fmt.Errorf("could not retrieve K8s Cluster maintenance window day")
		}

		tyme, ok := maintenanceWindow.GetTimeOk()
		if !ok || tyme == nil {
			return nil, fmt.Errorf("could not retrieve K8s Cluster maintenance window time")
		}

		temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)
	}

	return temp, nil
}

func ConvertK8sClustersToTable(clusters ionoscloud.KubernetesClusters) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve K8s Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertK8sClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertK8sNodepoolToTable(nodepool ionoscloud.KubernetesNodePool) ([]map[string]interface{}, error) {
	properties, ok := nodepool.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve K8s Nodepool properties")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.K8sNodepool, nodepool)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
	if ok && maintenanceWindow != nil {
		day, ok := maintenanceWindow.GetDayOfTheWeekOk()
		if !ok || day == nil {
			return nil, fmt.Errorf("could not retrieve K8s Nodepool maintenance window day")
		}

		tyme, ok := maintenanceWindow.GetTimeOk()
		if !ok || tyme == nil {
			return nil, fmt.Errorf("could not retrieve K8s Nodepool maintenance window time")
		}

		temp[0]["MaintenanceWindow"] = fmt.Sprintf("%s %s", *day, *tyme)
	}

	return temp, nil
}

func ConvertK8sNodepoolsToTable(nodepools ionoscloud.KubernetesNodePools) ([]map[string]interface{}, error) {
	items, ok := nodepools.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve K8s Nodepools items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertK8sNodepoolToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertRequestsToTable(requests ionoscloud.Requests) ([]map[string]interface{}, error) {
	items, ok := requests.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("failed to retrieve Requests items")
	}

	res := make([]map[string]interface{}, 0)
	for _, item := range *items {
		temp, err := ConvertRequestToTable(item)
		if err != nil {
			return nil, err
		}

		res = append(res, temp...)
	}

	return res, nil
}

func ConvertRequestToTable(request ionoscloud.Request) ([]map[string]interface{}, error) {
	metadata, ok := request.GetMetadataOk()
	if !ok || metadata == nil {
		return nil, fmt.Errorf("failed to retrieve Request metadata")
	}

	reqStatus, ok := metadata.GetRequestStatusOk()
	if !ok || reqStatus == nil {
		return nil, fmt.Errorf("failed to retrieve Request Status")
	}

	reqStatusMetadata, ok := reqStatus.GetMetadataOk()
	if !ok || reqStatusMetadata == nil {
		return nil, fmt.Errorf("failed to retrieve Request Status metadata")
	}

	targets, ok := reqStatusMetadata.GetTargetsOk()
	if !ok || targets == nil {
		return nil, fmt.Errorf("failed to retrieve Request Targets")
	}

	targetsInfo := make([]interface{}, 0)
	for _, target := range *targets {
		targetOk, ok := target.GetTargetOk()
		if !ok || targetOk == nil {
			continue
		}

		idOk, ok := targetOk.GetIdOk()
		if !ok || idOk == nil {
			continue
		}

		typeOk, ok := targetOk.GetTypeOk()
		if !ok || typeOk == nil {
			continue
		}

		targetsInfo = append(targetsInfo, fmt.Sprintf("%s (%s)", *idOk, string(*typeOk)))
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.Request, request)
	if err != nil {
		return nil, fmt.Errorf("failed to convert from JSON to Table format: %w", err)
	}

	temp[0]["Targets"] = targetsInfo

	return temp, nil
}
