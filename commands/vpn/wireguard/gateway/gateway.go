package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"ID", "Name", "PublicKey", "Description", "GatewayIP", "InterfaceIPv4", "InterfaceIPv6", "DatacenterId", "LanId", "ConnectionIPv4", "ConnectionIPv6", "InterfaceIP", "ListenPort", "Status"}
	// we can safely include both InterfaceIPv4 and InterfaceIPv6 cols because if the respective column empty, it won't be shown
	defaultCols = []string{"ID", "Name", "PublicKey", "Description", "GatewayIP", "InterfaceIPv4", "InterfaceIPv6", "DatacenterId", "ListenPort", "Status"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gateway",
			Short:            "Manage Wireguard VPN Gateways",
			Aliases:          []string{"g", "gw"},
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Get())
	cmd.AddCommand(Update())

	return cmd
}
