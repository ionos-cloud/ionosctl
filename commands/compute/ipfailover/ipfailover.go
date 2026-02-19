package ipfailover

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultIpFailoverCols = []string{"NicId", "Ip"}
)

func IpfailoverCmd() *core.Command {
	ipfailoverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipfailover",
			Aliases:          []string{"ipf"},
			Short:            "IP Failover Operations",
			Long:             "The sub-command of `ionosctl compute ipfailover` allows you to see information about IP Failovers groups available on a LAN, to add/remove IP Failover group from a LAN.",
			TraverseChildren: true,
		},
	}
	globalFlags := ipfailoverCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultIpFailoverCols, tabheaders.ColsMessage(defaultIpFailoverCols))
	_ = viper.BindPFlag(core.GetFlagName(ipfailoverCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = ipfailoverCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultIpFailoverCols, cobra.ShellCompDirectiveNoFileComp
	})

	ipfailoverCmd.AddCommand(IpFailoverListCmd())
	ipfailoverCmd.AddCommand(IpFailoverAddCmd())
	ipfailoverCmd.AddCommand(IpFailoverRemoveCmd())

	return core.WithConfigOverride(ipfailoverCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
