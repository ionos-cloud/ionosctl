package wireguard

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/peer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

/*
WireGuard is a modern VPN protocol that is lightweight, fast, and simple. It uses cryptographic primitives like Curve25519, ChaCha20, Poly1305, and BLAKE2s.
A WireGuard gateway establishes secure VPN connections using WireGuard peers (which act as VPN endpoints).
WireGuard operates at the IP layer (Layer 3) and uses a stateless model with fewer moving parts.
The configuration is static, meaning each peer (client or another gateway) is predefined in the configuration, which simplifies the overall setup.
*/

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "wireguard",
			Short:            "Manage Wireguard VPN Resources",
			Aliases:          []string{"wg"},
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(peer.Root())

	return cmd
}
