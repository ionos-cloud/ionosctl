package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gateway",
			Short:            "Manage IPsec VPN gateways",
			Aliases:          []string{"g"},
			TraverseChildren: true,
		},
	}
	// cmd.AddCommand()

	return cmd
}
