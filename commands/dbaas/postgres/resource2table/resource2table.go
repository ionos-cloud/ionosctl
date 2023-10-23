package resource2table

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

func ConvertAPIVersionToTable(apiVersion ionoscloud.APIVersion) ([]map[string]interface{}, error) {
	swaggerUrlOk, ok := apiVersion.GetSwaggerUrlOk()
	if !ok || swaggerUrlOk == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL API Version swagger URL")
	}

	if strings.HasPrefix(*swaggerUrlOk, "appserver:8181/postgresql") {
		*swaggerUrlOk = strings.TrimPrefix(*swaggerUrlOk, "appserver:8181/postgresql")
	}
	if !strings.HasPrefix(*swaggerUrlOk, ionoscloud.DefaultIonosServerUrl) {
		*swaggerUrlOk = fmt.Sprintf("%s%s", ionoscloud.DefaultIonosServerUrl, *swaggerUrlOk)
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.ApiVersion, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["SwaggerUrl"] = *swaggerUrlOk

	return temp, nil
}

func ConvertClusterToTable(cluster ionoscloud.ClusterResponse) ([]map[string]interface{}, error) {
	properties, ok := cluster.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster properties")
	}

	maintenanceWindow, ok := properties.GetMaintenanceWindowOk()
	if !ok || maintenanceWindow == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster maintenance window")
	}

	day, ok := maintenanceWindow.GetDayOfTheWeekOk()
	if !ok || day == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster maintenance window day")
	}

	tyme, ok := maintenanceWindow.GetTimeOk()
	if !ok || tyme == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster maintenance window time")
	}

	storage, ok := properties.GetStorageSizeOk()
	if !ok || storage == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster storage size")
	}

	ram, ok := properties.GetRamOk()
	if !ok || ram == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster RAM")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.Cluster, cluster)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["MaintenanceWindow"] = fmt.Sprintf("%v %v", *day, *tyme)
	temp[0]["RAM"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*ram), convbytes.MB, convbytes.GB))
	temp[0]["StorageSize"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*storage), convbytes.MB, convbytes.GB))

	connections, ok := properties.GetConnectionsOk()
	if ok && connections != nil {
		for _, con := range *connections {
			dcId, ok := con.GetDatacenterIdOk()
			if !ok || dcId == nil {
				return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster datacenter ID")
			}

			lanId, ok := con.GetLanIdOk()
			if !ok || lanId == nil {
				return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster lan ID")
			}

			cidr, ok := con.GetCidrOk()
			if !ok || cidr == nil {
				return nil, fmt.Errorf("could not retrieve PostgreSQL Cluster CIDRs")
			}

			temp[0]["DatacenterId"] = *dcId
			temp[0]["LanId"] = *lanId
			temp[0]["Cidr"] = *cidr
		}
	}

	return temp, nil
}

func ConvertClustersToTable(clusters ionoscloud.ClusterList) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertLogsToTable(logs *[]ionoscloud.ClusterLogsInstances) ([]map[string]interface{}, error) {
	if logs == nil {
		return nil, fmt.Errorf("no logs to process")
	}

	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		if instance.GetMessages() == nil {
			continue
		}

		for msgIdx, msg := range *instance.GetMessages() {
			o, err := json2table.ConvertJSONToTable("", jsonpaths.LogsMessage, msg)
			if err != nil {
				return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}

			o[0]["Instance"] = idx
			o[0]["MessageNumber"] = msgIdx
			if instance.GetName() != nil {
				o[0]["Name"] = *instance.GetName()
			}

			out = append(out, o...)
		}
	}

	return out, nil
}
