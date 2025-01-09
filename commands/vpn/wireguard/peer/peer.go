package peer

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

/*
A WireGuard peer is any device (client, server, or another gateway) that participates in a WireGuard VPN. Peers are identified by public/private key pairs.
WireGuard does not need complex negotiation (like IPsec IKE phases). Once two peers know each other’s public keys and IP addresses, they can connect instantly.
WireGuard is stateless: no persistent state is stored between connections, and packets are exchanged only when needed.
There is no session or tunnel establishment process like in IPsec. Instead, WireGuard peers exchange packets as needed without keeping an active session.
*/

var (
	allCols = []string{"ID", "Name", "Description", "Host", "Port", "WhitelistIPs", "PublicKey", "Status"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "peer",
			Short:            "Manage Wireguard VPN Peers",
			Aliases:          []string{"p"},
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Update())

	return cmd
}
