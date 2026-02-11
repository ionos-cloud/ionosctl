package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
)

func BackupsIds() []string {
	backupList, _, err := client.Must().PostgresClientV2.BackupsApi.BackupsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := backupList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range dataOk {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ids = append(ids, *itemId)
			}
		}
	} else {
		return nil
	}
	return ids
}

func BackupsIdsForCluster(clusterId string) []string {
	backupList, _, err := client.Must().PostgresClientV2.BackupsApi.BackupsGet(context.Background()).FilterClusterId(clusterId).Execute()
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := backupList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range dataOk {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ids = append(ids, *itemId)
			}
		}
	} else {
		return nil
	}
	return ids
}

func ClustersIds() []string {
	clusterList, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	convertedClusterList, err := json2table.ConvertJSONToTable(
		"items", jsonpaths.DbaasPostgresCluster, clusterList,
	)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(convertedClusterList, "ClusterId").AddInfo("DisplayName").AddInfo(
		"Location",
		"(%v)",
	).ToString()
}

func PostgresVersions() []string {
	versionList, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	versions := make([]string, 0)
	if dataOk, ok := versionList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range dataOk {
			completion := ""
			if item.Properties.Version != nil {
				completion += *item.Properties.Version
			}
			if item.Properties.Status != nil {
				completion += "\t status = " + *item.Properties.Status
			}
			if item.Properties.Comment != nil {
				completion += "; (" + *item.Properties.Comment + ")"
			}

			if completion != "" {
				versions = append(versions, completion)
			}
		}
	} else {
		return nil
	}
	return versions
}

func VersionsIds() []string {
	versionList, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := versionList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range dataOk {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ids = append(ids, *itemId)
			}
		}
	} else {
		return nil
	}
	return ids
}
