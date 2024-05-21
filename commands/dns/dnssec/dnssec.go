package dnssec

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"ID", "KeyTag", "DigestAlgorithmMnemonic", "Digest",
		"Flags", "PubKey", "ComposedKeyData", "Algorithm", "NsecMode"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dnssec",
			Aliases:          []string{"sec", "dnskey", "key", "keys"},
			Short:            "The sub-commands of `ionosctl dns dnssec` allow you to manage your DNSSEC Keys",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Get())

	return cmd
}
