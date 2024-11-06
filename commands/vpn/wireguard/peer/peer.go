package peer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
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

// PeersProperty returns a list of properties of all peers matching the given filters
func PeersProperty[V any](gatewayID string, f func(peer vpn.WireguardPeerRead) V, fs ...Filter) []V {
	recs, err := Peers(gatewayID, fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Peers returns all distributions matching the given filters
func Peers(gatewayID string, fs ...Filter) (vpn.WireguardPeerReadList, error) {
	req := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersGet(context.Background(), gatewayID)
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vpn.WireguardPeerReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vpn.WireguardPeerReadList{}, err
	}
	return ls, nil
}

type Filter func(request vpn.ApiWireguardgatewaysPeersGetRequest) (vpn.ApiWireguardgatewaysPeersGetRequest, error)
