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
			Use:              "k8s",
			Short:            "Kubernetes Operations",
			Long:             "The sub-commands of `ionosctl k8s` allow you to list, get, create, update, delete Kubernetes Clusters.",
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
