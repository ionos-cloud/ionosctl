package dnssec

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "KeyTag", "DigestAlgorithmMnemonic", "Digest", "Validity",
		"Flags", "PubKey", "ComposedKeyData", "Algorithm", "KskBits", "ZskBits", "NsecMode", "Nsec3Iterations", "Nsec3SaltBits"}

	defaultCols = []string{"Id", "KeyTag", "DigestAlgorithmMnemonic", "Digest", "Validity"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "dnssec",
			Aliases:          []string{"sec", "dnskey", "key", "keys"},
			Short:            "The sub-commands of 'ionosctl dns dnssec' allow you to manage your DNSSEC Keys",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, defaultCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Get())
	cmd.AddCommand(Create())
	cmd.AddCommand(Delete())

	return cmd
}
