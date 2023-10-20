package cluster

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Version", "MaintenanceWindow", "DatacenterId", "State"}
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Dataplatform Cluster Operations",
			Long:             "This command allows you to interact with the already created clusters or creates new clusters in your virtual data center",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	clusterCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClustersKubeConfigCmd())

	return clusterCmd
}
