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
			Use:   "lan",
			Short: "Kubernetes NodePool LAN Operations",
			Long: `Manage the LANs attached to a node pool's worker Nodes.

Attaching a LAN gives the pool's Nodes a network interface on that LAN (an
existing LAN in the same Data Center as the pool). Per LAN you may also define
routes - a destination network (CIDR) reachable via a gateway IP - so Nodes can
reach networks that sit behind that gateway. Use these sub-commands to list the
LANs on a pool, add a LAN (optionally with routes), or remove one.`,
			TraverseChildren: true,
		},
	}

	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanListCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanAddCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanRemoveCmd())

	return core.WithConfigOverride(k8sNodePoolLanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
