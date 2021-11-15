package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
)

func K8sClusterUpgradeVersions(outErr io.Writer, clusterId string) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.Background())
	cluster, _, err := k8sSvc.GetCluster(clusterId)
	clierror.CheckError(err, outErr)
	if cluster.Properties == nil || cluster.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *cluster.Properties.AvailableUpgradeVersions
}

func K8sNodePoolUpgradeVersions(outErr io.Writer, clusterId, nodepoolId string) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.Background())
	nodepool, _, err := k8sSvc.GetNodePool(clusterId, nodepoolId)
	clierror.CheckError(err, outErr)
	if nodepool.Properties == nil || nodepool.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *nodepool.Properties.AvailableUpgradeVersions
}
