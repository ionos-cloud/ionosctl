package node

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultK8sNodeCols = []string{"NodeId", "Name", "K8sVersion", "PublicIP", "PrivateIP", "State"}
)

func K8sNodeCmd() *core.Command {
	k8sNodeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "node",
			Aliases:          []string{"n"},
			Short:            "Kubernetes Node Operations",
			Long:             "The sub-commands of `ionosctl k8s node` allow you to list, get, recreate, delete Kubernetes Nodes.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sNodeCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultK8sNodeCols, tabheaders.ColsMessage(defaultK8sNodeCols))
	_ = viper.BindPFlag(core.GetFlagName(k8sNodeCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = k8sNodeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodeCols, cobra.ShellCompDirectiveNoFileComp
	})

	k8sNodeCmd.AddCommand(K8sNodeListCmd())
	k8sNodeCmd.AddCommand(K8sNodeGetCmd())
	k8sNodeCmd.AddCommand(K8sNodeRecreateCmd())
	k8sNodeCmd.AddCommand(K8sNodeDeleteCmd())

	return core.WithConfigOverride(k8sNodeCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
