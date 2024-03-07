package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/spf13/viper"
)

func BackupsIds() []string {
	clustersService := resources.NewBackupsService(client.Must(), context.Background())
	backupList, _, err := clustersService.List()
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := backupList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range *dataOk {
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
	clustersService := resources.NewBackupsService(client.Must(), context.Background())
	backupList, _, err := clustersService.ListBackups(clusterId)
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := backupList.GetItemsOk(); ok && dataOk != nil {
		for _, item := range *dataOk {
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
	clustersService := resources.NewClustersService(client.Must(), context.Background())
	clusterList, _, err := clustersService.List("")
	if err != nil {
		return nil
	}

	convertedClusterList, err := json2table.ConvertJSONToTable(
		"items", jsonpaths.DbaasPostgresCluster, clusterList.ClusterList,
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
	versionsService := resources.NewVersionsService(client.Must(), context.Background())
	versionList, _, err := versionsService.List()
	if err != nil {
		return nil
	}
	versions := make([]string, 0)
	if dataOk, ok := versionList.GetDataOk(); ok && dataOk != nil {
		for _, item := range *dataOk {
			if nameOk, ok := item.GetNameOk(); ok && nameOk != nil {
				versions = append(versions, *nameOk)
			}
		}
	} else {
		return nil
	}
	return versions
}

func UserNames(c *core.Command) []string {
	clusterId := viper.GetString(core.GetFlagName(c.NS, "cluster-id"))

	userList, _, err := client.Must().PostgresClient.UsersApi.UsersList(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}

	convertedUserList, err := json2table.ConvertJSONToTable(
		"items", jsonpaths.DbaasPostgresUser, userList,
	)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(convertedUserList, "Username").ToString()
}

func DatabaseNames(c *core.Command) []string {
	clusterId := viper.GetString(core.GetFlagName(c.NS, "cluster-id"))

	databaseList, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesList(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}

	convertedDatabaseList, err := json2table.ConvertJSONToTable(
		"items", jsonpaths.DbaasPostgresDatabase, databaseList,
	)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(convertedDatabaseList, "Name").ToString()
}
