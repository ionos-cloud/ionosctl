package ipsec

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/tunnel"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "ipsec",
			Short: "Manage IPSec VPN gateways and tunnels",
			Long: `Manage IPSec VPN resources.

IPSec is the standards-based site-to-site VPN — use it to interoperate with existing firewalls/routers. Unlike WireGuard, both ends negotiate the connection in two phases: IKE (phase 1) authenticates the peers and builds a secure control channel, and ESP (phase 2) is the negotiated channel that actually encrypts your traffic.

Model it as one gateway (the IONOS side on your LAN, IKEv2) plus one tunnel per remote site — the tunnel carries the remote host, the authentication (PSK/RSA), the IKE/ESP crypto parameters, and which local vs remote subnets may cross.`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(tunnel.Root())

	return cmd
}
