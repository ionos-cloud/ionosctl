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
			Use:              "node",
			Aliases:          []string{"n"},
			Short:            "Kubernetes Node Operations",
			Long:             "The sub-commands of `ionosctl compute k8s node` allow you to list, get, recreate, delete Kubernetes Nodes.",
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
