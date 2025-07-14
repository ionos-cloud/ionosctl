package completer

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func K8sClusterUpgradeVersions(clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	k8sSvc := resources.NewK8sService(client, context.Background())
	cluster, _, err := k8sSvc.GetCluster(clusterId, resources.QueryParams{})
	if err != nil {
		return nil
	}
	if cluster.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return cluster.Properties.AvailableUpgradeVersions
}

func K8sNodePoolUpgradeVersions(clusterId, nodepoolId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	k8sSvc := resources.NewK8sService(client, context.Background())
	nodepool, _, err := k8sSvc.GetNodePool(clusterId, nodepoolId, resources.QueryParams{})
	if err != nil {
		return nil
	}
	if nodepool.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return nodepool.Properties.AvailableUpgradeVersions
}
