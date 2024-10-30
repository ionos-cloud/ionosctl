package vpn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vpn",
			Short:            "VPN Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(wireguard.Root())

	return cmd
}