package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

// clusterCompleterCols defines columns needed for cluster ID autocompletion.
var clusterCompleterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id"},
	{Name: "DisplayName", JSONPath: "properties.displayName"},
	{Name: "Location", JSONPath: "properties.location"},
}

// userCompleterCols defines columns needed for user name autocompletion.
var userCompleterCols = []table.Column{
	{Name: "Username", JSONPath: "properties.username"},
}

// databaseCompleterCols defines columns needed for database name autocompletion.
var databaseCompleterCols = []table.Column{
	{Name: "Name", JSONPath: "properties.name"},
}

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

	t := table.New(clusterCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(clusterList); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "ClusterId").AddInfo("DisplayName").AddInfo(
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
	clusterId := viper.GetString(core.GetFlagName(c.NS, "cluster-id"))

	userList, _, err := client.Must().PostgresClient.UsersApi.UsersList(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}

	t := table.New(userCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(userList); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Username").ToString()
}

func DatabaseNames(c *core.Command) []string {
	clusterId := viper.GetString(core.GetFlagName(c.NS, "cluster-id"))

	databaseList, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesList(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}

	t := table.New(databaseCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(databaseList); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Name").ToString()
}
