package vpn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vpn",
			Short:            "The sub-commands of the 'vpn' resource help automate VPN Gateways resources",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(wireguard.Root())
	cmd.AddCommand(ipsec.Root())

	return cmd
}
