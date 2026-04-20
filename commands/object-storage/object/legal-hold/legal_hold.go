package legalhold

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagKey       = "key"
	flagKeyShort  = "k"
	flagVersionId = "version-id"
	flagStatus    = "status"
)

var allCols = []table.Column{
	{Name: "Status", JSONPath: "Status", Default: true},
}

type legalHoldInfo struct {
	Status string `json:"Status"`
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "legal-hold",
			Aliases:          []string{"lh"},
			Short:            "Manage Object Lock legal hold for objects",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(GetCmd())
	cmd.AddCommand(PutCmd())

	return cmd
}
