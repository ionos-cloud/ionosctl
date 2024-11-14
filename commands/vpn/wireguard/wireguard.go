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
