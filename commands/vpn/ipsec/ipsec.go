package ipsec

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/tunnel"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

/*
IPsec (Internet Protocol Security) is a more traditional and widely adopted VPN protocol. It provides encryption and integrity for IP packets using various cryptographic suites.
An IPsec gateway creates IPsec tunnels that encapsulate and secure traffic between two network endpoints.
IPsec operates in two main modes: Transport mode (securing traffic between two hosts) and Tunnel mode (securing traffic between two networks).
IPsec is more complex to set up compared to WireGuard, with phases of negotiation, such as IKE (Internet Key Exchange) for key management.
*/

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipsec",
			Short:            "Manage IPsec resources",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(tunnel.Root())

	return cmd
}
