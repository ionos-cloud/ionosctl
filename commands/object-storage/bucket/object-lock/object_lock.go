package objectlock

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagMode  = "mode"
	flagDays  = "days"
	flagYears = "years"
)

var allCols = []table.Column{
	{Name: "ObjectLockEnabled", JSONPath: "ObjectLockEnabled", Default: true},
	{Name: "RetentionMode", JSONPath: "RetentionMode", Default: true},
	{Name: "RetentionDays", JSONPath: "RetentionDays", Default: true},
	{Name: "RetentionYears", JSONPath: "RetentionYears"},
}

type configInfo struct {
	ObjectLockEnabled string `json:"ObjectLockEnabled"`
	RetentionMode     string `json:"RetentionMode"`
	RetentionDays     string `json:"RetentionDays"`
	RetentionYears    string `json:"RetentionYears"`
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "object-lock",
			Aliases:          []string{"ol"},
			Short:            "Manage bucket Object Lock configuration",
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
