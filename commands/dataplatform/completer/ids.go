package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
	"github.com/spf13/viper"
)

func ClustersIds(outErr io.Writer) []string {
	client, err := getDataPlatformClient()
	clierror.CheckError(err, outErr)
	clustersService := resources.NewClustersService(client, context.TODO())
	clusterList, _, err := clustersService.List("")
	clierror.CheckError(err, outErr)
	ids := make([]string, 0)
	if dataOk, ok := clusterList.ClusterListResponseData.GetItemsOk(); ok && dataOk != nil {
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

func NodePoolsIds(outErr io.Writer, clusterId string) []string {
	client, err := getDataPlatformClient()
	clierror.CheckError(err, outErr)
	nodePoolsService := resources.NewNodePoolsService(client, context.TODO())
	nodePoolsList, _, err := nodePoolsService.List(clusterId)
	clierror.CheckError(err, outErr)
	ids := make([]string, 0)
	if dataOk, ok := nodePoolsList.NodePoolListResponseData.GetItemsOk(); ok && dataOk != nil {
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

func DataPlatformVersions(outErr io.Writer) []string {
	client, err := getDataPlatformClient()
	clierror.CheckError(err, outErr)
	versionsService := resources.NewVersionsService(client, context.TODO())
	versionsList, _, err := versionsService.List()
	clierror.CheckError(err, outErr)
	versions := make([]string, 0)
	if len(versionsList) > 0 {
		for _, item := range versionsList {

			versions = append(versions, item)

		}
	} else {
		return nil
	}
	return versions
}

// Get Client for Completion Functions
func getDataPlatformClient() (*resources.Client, error) {
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
