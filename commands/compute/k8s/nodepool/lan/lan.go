package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allK8sNodePoolLanCols = []table.Column{
	{Name: "LanId", JSONPath: "id", Default: true},
	{Name: "Dhcp", JSONPath: "dhcp", Default: true},
	{Name: "RoutesNetwork", JSONPath: "routes.*.network", Default: true},
	{Name: "RoutesGatewayIp", JSONPath: "routes.*.gatewayIp", Default: true},
}

func K8sNodePoolLanCmd() *core.Command {
	k8sNodePoolLanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "Kubernetes NodePool LAN Operations",
			Long:             "The sub-commands of `ionosctl compute k8s nodepool lan` allow you to list, add, remove Kubernetes Node Pool LANs.",
			TraverseChildren: true,
		},
	}

	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanListCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanAddCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanRemoveCmd())

	return core.WithConfigOverride(k8sNodePoolLanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
