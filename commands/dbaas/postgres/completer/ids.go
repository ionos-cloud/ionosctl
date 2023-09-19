package completer

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
)

func BackupsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	clustersService := resources.NewBackupsService(client, context.TODO())
	backupList, _, err := clustersService.List()
	if err != nil {
		die.Die(err.Error())
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

func BackupsIdsForCluster(_ io.Writer, clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	clustersService := resources.NewBackupsService(client, context.TODO())
	backupList, _, err := clustersService.ListBackups(clusterId)
	if err != nil {
		die.Die(err.Error())
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

func ClustersIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	clustersService := resources.NewClustersService(client, context.TODO())
	clusterList, _, err := clustersService.List("")
	if err != nil {
		die.Die(err.Error())
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

func PostgresVersions(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	versionsService := resources.NewVersionsService(client, context.TODO())
	versionList, _, err := versionsService.List()
	if err != nil {
		die.Die(err.Error())
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
