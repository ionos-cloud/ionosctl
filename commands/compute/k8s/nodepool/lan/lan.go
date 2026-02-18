package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultK8sNodePoolLanCols = []string{"LanId", "Dhcp", "RoutesNetwork", "RoutesGatewayIp"}
)

func K8sNodePoolLanCmd() *core.Command {
	k8sNodePoolLanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Short:            "Kubernetes NodePool LAN Operations",
			Long:             "The sub-commands of `ionosctl k8s nodepool lan` allow you to list, add, remove Kubernetes Node Pool LANs.",
			TraverseChildren: true,
		},
	}

	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanListCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanAddCmd())
	k8sNodePoolLanCmd.AddCommand(K8sNodePoolLanRemoveCmd())

	return core.WithConfigOverride(k8sNodePoolLanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
