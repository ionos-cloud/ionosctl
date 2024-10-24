package tunnel

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

/*
IPsec tunnels are used to establish secure communication paths between two endpoints. These endpoints could be between two gateways (gateway-to-gateway) or a gateway and a remote client.
IPsec requires a series of negotiations between the endpoints (IKE Phase 1 and 2) to establish the tunnel and negotiate encryption keys.
IPsec maintains a stateful session; meaning, it keeps track of active sessions and renegotiates security parameters when needed.
IPsec can handle dynamic negotiation of tunnel parameters, which can be more flexible in some use cases.
*/

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "tunnel",
			Short:            "Manage IPsec tunnels",
			Aliases:          []string{"t"},
			TraverseChildren: true,
		},
	}
	// cmd.AddCommand()

	return cmd
}
