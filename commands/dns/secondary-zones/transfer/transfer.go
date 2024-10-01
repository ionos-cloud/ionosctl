package transfer

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var allCols = []string{"PrimaryIP", "Status", "ErrorMessage"}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "transfer",
			Aliases:          []string{"t"},
			Short:            "",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, allCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(getCmd())
	cmd.AddCommand(startCmd())

	return cmd
}
