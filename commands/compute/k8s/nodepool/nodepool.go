package nodepool

import (
	nplan "github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/nodepool/lan"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultK8sNodePoolCols = []string{"NodePoolId", "Name", "K8sVersion", "NodeCount", "DatacenterId", "State"}
	allK8sNodePoolCols     = []string{
		"NodePoolId", "Name", "K8sVersion", "DatacenterId", "NodeCount", "CpuFamily", "ServerType", "StorageType", "State", "LanIds",
		"CoresCount", "RamSize", "AvailabilityZone", "StorageSize", "MaintenanceWindow", "AutoScaling", "PublicIps", "AvailableUpgradeVersions",
		"Annotations", "Labels", "ClusterId",
	}
)

func K8sNodePoolCmd() *core.Command {
	k8sNodePoolCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Kubernetes NodePool Operations",
			Long:             "The sub-commands of `ionosctl k8s nodepool` allow you to list, get, create, update, delete Kubernetes NodePools.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sNodePoolCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultK8sNodePoolCols, tabheaders.ColsMessage(allK8sNodePoolCols))
	_ = viper.BindPFlag(core.GetFlagName(k8sNodePoolCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = k8sNodePoolCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sNodePoolCols, cobra.ShellCompDirectiveNoFileComp
	})

	k8sNodePoolCmd.AddCommand(K8sNodePoolListCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolGetCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolCreateCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolUpdateCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolDeleteCmd())
	k8sNodePoolCmd.AddCommand(nplan.K8sNodePoolLanCmd())

	return core.WithConfigOverride(k8sNodePoolCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
