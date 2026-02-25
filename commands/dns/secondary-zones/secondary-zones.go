package secondary_zones

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/secondary-zones/transfer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.zoneName", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "PrimaryIPs", JSONPath: "properties.primaryIps", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "secondary-zone",
			Aliases:          []string{"secondary-zones", "sz"},
			Short:            "All commands related to secondary zones",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, table.AllCols(allCols), table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(transfer.Root())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(updateCmd())

	return cmd
}
