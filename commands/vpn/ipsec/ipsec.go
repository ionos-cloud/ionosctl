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
			Use:              "ipsec",
			Short:            "Manage ipsec VPN Resources",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(gateway.Root())
	cmd.AddCommand(tunnel.Root())

	return cmd
}
