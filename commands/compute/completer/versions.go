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
	cluster, _, err := k8sSvc.GetCluster(clusterId)
	if err != nil {
		return nil
	}
	if cluster.Properties == nil || cluster.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *cluster.Properties.AvailableUpgradeVersions
}

func K8sNodePoolUpgradeVersions(clusterId, nodepoolId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	k8sSvc := resources.NewK8sService(client, context.Background())
	nodepool, _, err := k8sSvc.GetNodePool(clusterId, nodepoolId)
	if err != nil {
		return nil
	}
	if nodepool.Properties == nil || nodepool.Properties.AvailableUpgradeVersions == nil {
		return nil
	}
	return *nodepool.Properties.AvailableUpgradeVersions
}

// K8sNodePoolTaints returns the node pool's taints as KEY[=VALUE]:EFFECT.
func K8sNodePoolTaints(clusterId, nodepoolId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	k8sSvc := resources.NewK8sService(client, context.Background())
	nodepool, _, err := k8sSvc.GetNodePool(clusterId, nodepoolId)
	if err != nil || nodepool.Properties == nil || nodepool.Properties.Taints == nil {
		return nil
	}
	out := make([]string, 0, len(*nodepool.Properties.Taints))
	for _, t := range *nodepool.Properties.Taints {
		if t.Key == nil || t.Effect == nil {
			continue
		}
		s := *t.Key
		if t.Value != nil && *t.Value != "" {
			s += "=" + *t.Value
		}
		out = append(out, s+":"+string(*t.Effect))
	}
	return out
}
