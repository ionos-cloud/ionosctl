package resource2table

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	sdkmongo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	sdkpsql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

func ConvertDbaasMongoClusterToTable(cluster sdkmongo.ClusterResponse) ([]map[string]interface{}, error) {
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

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasMongoCluster, cluster)
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

func ConvertDbaasMongoClustersToTable(clusters sdkmongo.ClusterList) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertDbaasMongoClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertMongoDbaasLogsToTable(logs *[]sdkmongo.ClusterLogsInstances) ([]map[string]interface{}, error) {
	if logs == nil {
		return nil, fmt.Errorf("no logs to process")
	}

	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		if instance.GetMessages() == nil {
			continue
		}
		for msgIdx, msg := range *instance.GetMessages() {
			o, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasLogsMessage, msg)
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

func ConvertDbaasMongoTemplateToTable(template sdkmongo.TemplateResponse) ([]map[string]interface{}, error) {
	properties, ok := template.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template properties")
	}

	ram, ok := properties.GetRamOk()
	if !ok || ram == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template RAM")
	}

	storage, ok := properties.GetStorageSizeOk()
	if !ok || storage == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Template storage")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasMongoTemplates, template)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["RAM"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*ram), convbytes.MB, convbytes.GB))
	temp[0]["StorageSize"] = fmt.Sprintf("%d GB", convbytes.Convert(int64(*storage), convbytes.MB, convbytes.GB))

	return temp, nil
}

func ConvertDbaasMongoTemplatesToTable(templates sdkmongo.TemplateList) ([]map[string]interface{}, error) {
	items, ok := templates.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Templates items")
	}

	var templatesConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertDbaasMongoTemplateToTable(item)
		if err != nil {
			return nil, err
		}

		templatesConverted = append(templatesConverted, temp...)
	}

	return templatesConverted, nil
}

func ConvertDbaasMongoUserToTable(user sdkmongo.User) ([]map[string]interface{}, error) {
	properties, ok := user.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Mongo User properties")
	}

	roles, ok := properties.GetRolesOk()
	if !ok || roles == nil {
		return nil, fmt.Errorf("could not retrieve Mongo User roles")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasMongoUser, user)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["Roles"] = strings.Join(functional.Map(*properties.GetRoles(), func(role sdkmongo.UserRoles) string {
		val, ok := role.GetRoleOk()
		if !ok {
			return ""
		}
		db, ok := role.GetDatabaseOk()
		if !ok {
			return ""
		}
		return fmt.Sprintf("%s: %s", *db, *val)
	}), ", ")

	return temp, nil
}

func ConvertDbaasMongoUsersToTable(users sdkmongo.UsersList) ([]map[string]interface{}, error) {
	items, ok := users.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Users items")
	}

	var usersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertDbaasMongoUserToTable(item)
		if err != nil {
			return nil, err
		}

		usersConverted = append(usersConverted, temp...)
	}

	return usersConverted, nil
}

func ConvertDbaasPostgresAPIVersionToTable(apiVersion sdkpsql.APIVersion) ([]map[string]interface{}, error) {
	swaggerUrlOk, ok := apiVersion.GetSwaggerUrlOk()
	if !ok || swaggerUrlOk == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL API Version swagger URL")
	}

	if strings.HasPrefix(*swaggerUrlOk, "appserver:8181/postgresql") {
		*swaggerUrlOk = strings.TrimPrefix(*swaggerUrlOk, "appserver:8181/postgresql")
	}
	if !strings.HasPrefix(*swaggerUrlOk, sdkpsql.DefaultIonosServerUrl) {
		*swaggerUrlOk = fmt.Sprintf("%s%s", sdkpsql.DefaultIonosServerUrl, *swaggerUrlOk)
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasPostgresApiVersion, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["SwaggerUrl"] = *swaggerUrlOk

	return temp, nil
}

func ConvertDbaasPostgresClusterToTable(cluster sdkpsql.ClusterResponse) ([]map[string]interface{}, error) {
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

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasPostgresCluster, cluster)
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

func ConvertDbaasPostgresClustersToTable(clusters sdkpsql.ClusterList) ([]map[string]interface{}, error) {
	items, ok := clusters.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL Clusters items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertDbaasPostgresClusterToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}

func ConvertDbaasPostgresLogsToTable(logs *[]sdkpsql.ClusterLogsInstances) ([]map[string]interface{}, error) {
	if logs == nil {
		return nil, fmt.Errorf("no logs to process")
	}

	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		if instance.GetMessages() == nil {
			continue
		}

		for msgIdx, msg := range *instance.GetMessages() {
			o, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasLogsMessage, msg)
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
