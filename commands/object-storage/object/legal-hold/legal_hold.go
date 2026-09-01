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
			Use:     "legal-hold",
			Aliases: []string{"lh"},
			Short:   "Manage Object Lock legal hold on objects",
			Long: `Manage the Object Lock legal hold on an object.

A legal hold is an on/off WORM protection: while it is ON, the object (version) cannot be deleted or overwritten. Unlike retention, it has NO expiry date - it stays in force until someone explicitly sets it OFF, and it can never be bypassed with --bypass-governance-retention.

Legal hold and retention are independent controls: an object can have both, either, or neither, and it stays locked while ANY of them is active. Typical use is preserving evidence for the duration of litigation, where the end date is unknown.

Set per object VERSION (pass --version-id to target a specific version). Requires a bucket created with Object Lock enabled; Object Lock cannot be added to an existing bucket and forces versioning on.`,
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
