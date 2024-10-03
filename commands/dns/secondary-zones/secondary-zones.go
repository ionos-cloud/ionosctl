package secondary_zones

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/secondary-zones/transfer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var allCols = []string{"Id", "Name", "Description", "PrimaryIPs", "State"}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "secondary-zone",
			Aliases:          []string{"secondary-zones", "sz"},
			Short:            "All commands related to secondary zones",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, allCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
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
