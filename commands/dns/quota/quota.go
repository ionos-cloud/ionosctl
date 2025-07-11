package quota

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"ZonesUsed", "ZonesLimit",
		"SecondaryZonesUsed", "SecondaryZonesLimit",
		"RecordsUsed", "RecordsLimit",
		"ReverseRecordsUsed", "ReverseRecordsLimit"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "quota",
			Aliases:          []string{"q"},
			Short:            "The sub-commands of 'ionosctl dns quota' allow you to see your DNS Quotas",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Get())

	return cmd
}
