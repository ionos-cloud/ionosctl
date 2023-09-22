package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
)

func BackupsIds(outErr io.Writer) []string {
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

func ClustersIds(outErr io.Writer) []string {
	clustersService := resources.NewClustersService(client.Must(), context.Background())
	clusterList, _, err := clustersService.List("")
	if err != nil {
		return nil
	}
	ids := make([]string, 0)
	if dataOk, ok := clusterList.ClusterList.GetItemsOk(); ok && dataOk != nil {
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

func PostgresVersions(outErr io.Writer) []string {
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
