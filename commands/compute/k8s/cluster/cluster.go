package cluster

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow", "Public", "Location"}
	allK8sClusterCols     = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow", "Public", "Location", "NatGatewayIp", "NodeSubnet", "AvailableUpgradeVersions", "ViableNodePoolVersions", "S3Bucket", "ApiSubnetAllowList"}
)

func K8sClusterCmd() *core.Command {
	k8sClusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Kubernetes Cluster Operations",
			Long:             "The sub-commands of `ionosctl compute k8s` allow you to perform Kubernetes Operations.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sClusterCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultK8sClusterCols, tabheaders.ColsMessage(allK8sClusterCols))
	_ = viper.BindPFlag(core.GetFlagName(k8sClusterCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = k8sClusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	k8sClusterCmd.AddCommand(K8sClusterListCmd())
	k8sClusterCmd.AddCommand(K8sClusterGetCmd())
	k8sClusterCmd.AddCommand(K8sClusterCreateCmd())
	k8sClusterCmd.AddCommand(K8sClusterUpdateCmd())
	k8sClusterCmd.AddCommand(K8sClusterDeleteCmd())

	return core.WithConfigOverride(k8sClusterCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
