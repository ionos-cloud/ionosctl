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
			Use:     "object-lock",
			Aliases: []string{"ol"},
			Short:   "Manage WORM (Write-Once-Read-Many) retention on a bucket",
			Long: `Manage a bucket's Object Lock configuration. Object Lock implements WORM (Write-Once-Read-Many): once written, an object version is protected from being overwritten or deleted until its retention period expires. This is used to meet regulatory/compliance requirements and to protect against accidental or malicious deletion (including ransomware).

Two retention modes exist:
  GOVERNANCE  Protects the object, but users with the special bypass permission can shorten or remove the retention. Use for internal data-protection policies.
  COMPLIANCE  Stricter: NOBODY - not even the account root - can shorten, override or delete the locked version until the retention period elapses. Use only when a legal/regulatory hold demands it, as it is genuinely irreversible.

Prerequisites and interplay: Object Lock can only be used on a bucket that was created with --object-lock (it cannot be enabled on an existing bucket), and it inherently requires versioning (enabling Object Lock enables versioning). The default retention set here applies to newly uploaded object versions; per-object overrides are set at upload time.`,
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
