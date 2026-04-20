package retention

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagKey                       = "key"
	flagKeyShort                  = "k"
	flagVersionId                 = "version-id"
	flagMode                      = "mode"
	flagRetainUntilDate           = "retain-until-date"
	flagBypassGovernanceRetention = "bypass-governance-retention"
)

var allCols = []table.Column{
	{Name: "Mode", JSONPath: "Mode", Default: true},
	{Name: "RetainUntilDate", JSONPath: "RetainUntilDate", Default: true},
}

type retentionInfo struct {
	Mode            string `json:"Mode"`
	RetainUntilDate string `json:"RetainUntilDate"`
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "retention",
			Aliases:          []string{"ret"},
			Short:            "Manage Object Lock retention for objects",
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
