package completer

import (
	"context"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func K8sClusterUpgradeVersions(_ io.Writer, clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.Background())
	cluster, _, err := k8sSvc.GetCluster(clusterId, resources.QueryParams{})
	if err != nil {
		die.Die(err.Error())
	}

	if cluster.Properties == nil || cluster.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *cluster.Properties.AvailableUpgradeVersions
}

func K8sNodePoolUpgradeVersions(_ io.Writer, clusterId, nodepoolId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.Background())
	nodepool, _, err := k8sSvc.GetNodePool(clusterId, nodepoolId, resources.QueryParams{})
	if err != nil {
		die.Die(err.Error())
	}

	if nodepool.Properties == nil || nodepool.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *nodepool.Properties.AvailableUpgradeVersions
}
