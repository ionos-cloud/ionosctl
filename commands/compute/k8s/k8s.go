package k8s

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/kubeconfig"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/node"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/nodepool"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func K8sCmd() *core.Command {
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:   "k8s",
			Short: "Managed Kubernetes (MK8s) Operations",
			Long: `Manage IONOS Managed Kubernetes (MK8s).

A Managed Kubernetes deployment has two layers: the managed control plane (API server, scheduler, etcd) that IONOS operates for you - you choose its Kubernetes version, a weekly maintenance window and optional API-server access controls - and one or more node pools, each a group of worker Nodes with identical hardware (cores, RAM, storage, CPU family) pinned to a single Data Center, where your workloads run.

Version skew: a node pool's Kubernetes version must be equal to or lower than the cluster's version. Upgrade the control plane first, then the node pools.`,
			TraverseChildren: true,
		},
	}

	k8sCmd.AddCommand(cluster.K8sClusterCmd())
	k8sCmd.AddCommand(nodepool.K8sNodePoolCmd())
	k8sCmd.AddCommand(node.K8sNodeCmd())
	k8sCmd.AddCommand(kubeconfig.K8sKubeconfigCmd())
	k8sCmd.AddCommand(version.K8sVersionCmd())

	return core.WithConfigOverride(k8sCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
