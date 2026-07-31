package node

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allK8sNodeCols = []table.Column{
	{Name: "NodeId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "K8sVersion", JSONPath: "properties.k8sVersion", Default: true},
	{Name: "PublicIP", JSONPath: "properties.publicIP", Default: true},
	{Name: "PrivateIP", JSONPath: "properties.privateIP", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func K8sNodeCmd() *core.Command {
	k8sNodeCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "node",
			Aliases: []string{"n"},
			Short:   "Kubernetes Node (individual worker) Operations",
			Long: `Operate on individual worker Nodes inside a node pool.

A Node is one worker machine created from its node pool's template. You do not
create Nodes directly - the pool's node count (or autoscaling) determines how
many exist. These sub-commands let you inspect a Node, recreate a faulty one, or
delete one. Nodes are always addressed by their cluster and node pool:
--cluster-id and --nodepool-id are required alongside --node-id.`,
			TraverseChildren: true,
		},
	}
	k8sNodeCmd.AddColsFlag(allK8sNodeCols)

	k8sNodeCmd.AddCommand(K8sNodeListCmd())
	k8sNodeCmd.AddCommand(K8sNodeGetCmd())
	k8sNodeCmd.AddCommand(K8sNodeRecreateCmd())
	k8sNodeCmd.AddCommand(K8sNodeDeleteCmd())

	return core.WithConfigOverride(k8sNodeCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
