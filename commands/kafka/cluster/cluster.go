package cluster

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Version", "Size", "DatacenterId", "LanId", "BrokerAddresses", "State", "StateMessage"}
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Short:            "The sub-commands of 'ionosctl kafka cluster' allow you to manage kafka clusters",
			Aliases:          []string{"cl"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(FindByID())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Create())
	return cmd
}
