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
			Use:     "retention",
			Aliases: []string{"ret"},
			Short:   "Manage Object Lock WORM retention on objects",
			Long: `Manage the Object Lock retention placed on an object.

Retention is a WORM (write-once-read-many) protection: it forbids deleting or overwriting an object (version) until a "retain-until" date passes, after which the protection lapses automatically. There are two modes:
  - GOVERNANCE: users holding the bypass permission can shorten or remove the lock (and delete with --bypass-governance-retention). Use for accidental-deletion protection with an escape hatch.
  - COMPLIANCE: nobody, not even the account root, can shorten or remove the lock before the date - used for regulatory mandates.

Retention is set per object VERSION; pass --version-id to target a version other than the current one. It differs from a legal hold, which has no expiry date and stays on until explicitly turned OFF; an object stays locked while EITHER is in force.

Requires a bucket created with Object Lock enabled. Object Lock cannot be added to an existing bucket and it forces versioning on permanently.`,
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
