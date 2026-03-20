package ipblock

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultIpBlockCols = []string{"IpBlockId", "Name", "Location", "Size", "Ips", "State"}
)

func IpblockCmd() *core.Command {
	ipblockCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Aliases:          []string{"ip", "ipb"},
			Short:            "IpBlock Operations",
			Long:             "The sub-commands of `ionosctl compute ipblock` allow you to create/reserve, list, get, update, delete IpBlocks.",
			TraverseChildren: true,
		},
	}
	globalFlags := ipblockCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultIpBlockCols, tabheaders.ColsMessage(defaultIpBlockCols))
	_ = viper.BindPFlag(core.GetFlagName(ipblockCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = ipblockCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultIpBlockCols, cobra.ShellCompDirectiveNoFileComp
	})

	ipblockCmd.AddCommand(IpBlockListCmd())
	ipblockCmd.AddCommand(IpBlockGetCmd())
	ipblockCmd.AddCommand(IpBlockCreateCmd())
	ipblockCmd.AddCommand(IpBlockUpdateCmd())
	ipblockCmd.AddCommand(IpBlockDeleteCmd())

	return core.WithConfigOverride(ipblockCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
