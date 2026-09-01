package wireguard

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/peer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "wireguard",
			Short: "Manage WireGuard VPN gateways and peers",
			Long: `Manage WireGuard VPN resources.

WireGuard connects two endpoints that each hold a key pair and know the other's PUBLIC key — there is no session handshake or phase negotiation like IPSec, so tunnels come up instantly and statelessly.

Model it as: create one 'gateway' (the IONOS side: a key pair, a tunnel interface IP and a UDP listen port on your LAN), then add one 'peer' per remote device (its public key, the IPs it is allowed to send, and where to reach it). Traffic flows once both sides list each other.`,
			Aliases:          []string{"wg"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(peer.Root())

	return cmd
}
