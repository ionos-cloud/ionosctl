package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
)

func BackupsIds() []string {
	backupList, _, err := client.Must().PostgresClient.BackupsApi.ClustersBackupsGet(context.Background()).Execute()
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
	backupList, _, err := client.Must().PostgresClient.BackupsApi.
		ClusterBackupsGet(context.Background(), clusterId).Execute()
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
	clusterList, _, err := client.Must().PostgresClient.ClustersApi.ClustersGet(context.Background()).Execute()
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
	versionList, _, err := client.Must().PostgresClient.ClustersApi.PostgresVersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	versions := make([]string, 0)
	if dataOk, ok := versionList.GetDataOk(); ok && dataOk != nil {
		for _, item := range dataOk {
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
	clusterId, err := c.Command.Flags().GetString("cluster-id")
	if err != nil {
		return nil
	}

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
	clusterId, err := c.Command.Flags().GetString("cluster-id")
	if err != nil {
		return nil
	}

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
