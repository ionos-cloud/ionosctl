package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	"github.com/spf13/viper"
)

func BackupsIds(outErr io.Writer) []string {
	client, err := getDbaasPgsqlClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewBackupsService(client, context.TODO())
	backupList, _, err := clustersService.List()
	clierror.CheckError(err, outErr)
	ids := make([]string, 0)
	if dataOk, ok := backupList.GetDataOk(); ok && dataOk != nil {
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
	client, err := getDbaasPgsqlClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewClustersService(client, context.TODO())
	clusterList, _, err := clustersService.List("")
	clierror.CheckError(err, outErr)
	ids := make([]string, 0)
	if dataOk, ok := clusterList.ClusterList.GetDataOk(); ok && dataOk != nil {
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
	client, err := getDbaasPgsqlClient()
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

// Get Client for Completion Functions
func getDbaasPgsqlClient() (*resources.Client, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}
