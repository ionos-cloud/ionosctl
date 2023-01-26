package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
)

func BackupsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewBackupsService(client, context.TODO())
	backupList, _, err := clustersService.List()
	clierror.CheckError(err, outErr)
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

func BackupsIdsForCluster(outErr io.Writer, clusterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewBackupsService(client, context.TODO())
	backupList, _, err := clustersService.ListBackups(clusterId)
	clierror.CheckError(err, outErr)
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
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewClustersService(client, context.TODO())
	clusterList, _, err := clustersService.List("")
	clierror.CheckError(err, outErr)
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
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	versionsService := resources.NewVersionsService(client, context.TODO())
	versionList, _, err := versionsService.List()
	clierror.CheckError(err, outErr)
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
