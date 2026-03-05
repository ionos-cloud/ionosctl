package quota

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "ZonesUsed", JSONPath: "quotaUsage.zones", Default: true},
	{Name: "ZonesLimit", JSONPath: "quotaLimits.zones", Default: true},
	{Name: "SecondaryZonesUsed", JSONPath: "quotaUsage.secondaryZones", Default: true},
	{Name: "SecondaryZonesLimit", JSONPath: "quotaLimits.secondaryZones", Default: true},
	{Name: "RecordsUsed", JSONPath: "quotaUsage.records", Default: true},
	{Name: "RecordsLimit", JSONPath: "quotaLimits.records", Default: true},
	{Name: "ReverseRecordsUsed", JSONPath: "quotaUsage.reverseRecords", Default: true},
	{Name: "ReverseRecordsLimit", JSONPath: "quotaLimits.reverseRecords", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "quota",
			Aliases:          []string{"q"},
			Short:            "The sub-commands of 'ionosctl dns quota' allow you to see your DNS Quotas",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Get())

	return cmd
}
