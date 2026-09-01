package encryption

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "SSEAlgorithm", JSONPath: "SSEAlgorithm", Default: true},
}

type encryptionRuleInfo struct {
	SSEAlgorithm string `json:"SSEAlgorithm"`
}

func EncryptionCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "encryption",
			Aliases: []string{"enc"},
			Short:   "Manage default server-side encryption (SSE) for a bucket",
			Long: `Manage a bucket's default server-side encryption (SSE). When a default encryption rule is set, every new object written to the bucket is encrypted at rest by the server unless the upload request specifies its own encryption. This is transparent to readers: objects are decrypted on the fly for authorized GET requests, so no client-side key handling is required.

Default encryption applies only to objects created AFTER the rule is set; it does not retroactively encrypt existing objects. 'put' sets/replaces the rule, 'get' shows the configured SSE algorithm, and 'delete' removes the default (new objects are then stored without automatic server-side encryption unless the upload requests it).`,
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
	cmd.AddCommand(DeleteCmd())

	return cmd
}
