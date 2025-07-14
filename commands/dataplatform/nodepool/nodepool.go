package nodepool

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

func NodepoolCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Dataplatform Nodepool Operations",
			Long:             "Node pools are the resources that powers the DataPlatformCluster",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(NodepoolListCmd())
	cmd.AddCommand(NodepoolCreateCmd())
	cmd.AddCommand(NodepoolGetCmd())
	cmd.AddCommand(NodepoolUpdateCmd())
	cmd.AddCommand(NodepoolDeleteCmd())

	return cmd
}

var (
	allCols = []string{"Id", "Name", "Nodes", "Cores", "CpuFamily", "Ram", "Storage", "MaintenanceWindow", "State",
		"AvailabilityZone", "Labels", "Annotations", "ClusterId"}
	defaultCols = []string{"Id", "Name", "Nodes", "Cores", "CpuFamily", "Ram", "Storage", "MaintenanceWindow", "State"}
)
