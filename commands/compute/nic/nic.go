package nic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultNicCols = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "IPv6Ips", "State"}
	allNicCols     = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "IPv6Ips", "State", "FirewallActive",
		"FirewallType", "DeviceNumber", "PciSlot", "Mac", "DHCPv6", "IPv6CidrBlock"}
)

func NicCmd() *core.Command {
	nicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Network Interfaces Operations",
			Long:             "The sub-commands of `ionosctl compute nic` allow you to create, list, get, update, delete NICs. To attach a NIC to a Load Balancer, use the Load Balancer command `ionosctl compute loadbalancer nic attach`.",
			TraverseChildren: true,
		},
	}
	globalFlags := nicCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNicCols, tabheaders.ColsMessage(allNicCols))
	_ = viper.BindPFlag(core.GetFlagName(nicCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = nicCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})

	nicCmd.AddCommand(NicListCmd())
	nicCmd.AddCommand(NicGetCmd())
	nicCmd.AddCommand(NicCreateCmd())
	nicCmd.AddCommand(NicUpdateCmd())
	nicCmd.AddCommand(NicDeleteCmd())

	return core.WithConfigOverride(nicCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
