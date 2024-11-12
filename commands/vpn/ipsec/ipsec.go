package ipsec

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/tunnel"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

/*
IPSec is a modern VPN protocol that is lightweight, fast, and simple. It uses cryptographic primitives like Curve25519, ChaCha20, Poly1305, and BLAKE2s.
A IPSec gateway establishes secure VPN connections using IPSec tunnels (which act as VPN endpoints).
IPSec operates at the IP layer (Layer 3) and uses a stateless model with fewer moving parts.
The configuration is static, meaning each tunnel (client or another gateway) is predefined in the configuration, which simplifies the overall setup.
*/

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipsec",
			Short:            "Manage ipsec VPN Resources",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(tunnel.Root())

	return cmd
}
