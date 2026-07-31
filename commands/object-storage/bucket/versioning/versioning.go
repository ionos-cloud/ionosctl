package versioning

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "Versioning", JSONPath: "Versioning", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "versioning",
			Aliases: []string{"ver"},
			Short:   "Manage object versioning on a bucket (keep multiple versions)",
			Long: `Manage a bucket's versioning state. With versioning Enabled, each overwrite or delete of a key keeps the previous version instead of discarding it, so you can recover from accidental overwrites and deletions (a "delete" just writes a delete marker; the data remains until you delete that specific version).

States: a fresh bucket is unversioned ("Disabled" here). You can set it to "Enabled" or "Suspended". Versioning cannot be turned fully off once used - Suspended stops creating NEW versions but keeps all versions already stored. Versioning is also a prerequisite for Object Lock/WORM.

Lifecycle interplay: on a versioned bucket, deleting current objects (or a lifecycle Expiration rule) only adds delete markers and leaves prior versions - and their storage cost - in place. Reclaim them with a lifecycle NoncurrentVersionExpiration rule. Use 'get' to read the current state and 'set' to change it.`,
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(GetCmd())
	cmd.AddCommand(SetCmd())

	return cmd
}
