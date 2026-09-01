package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "NatGatewayLanId", JSONPath: "id", Default: true},
	{Name: "GatewayIps", JSONPath: "gatewayIps", Default: true},
}

func NatgatewayLanCmd() *core.Command {
	natgatewayLanCmd := &core.Command{
		Command: &cobra.Command{
			Use:   "lan",
			Short: "NAT Gateway Lan Operations",
			Long: `These sub-commands manage which private LANs a NAT Gateway is attached to. Attaching a LAN is what lets servers on that LAN route their outbound traffic through the gateway: the gateway becomes reachable on the LAN via one or more gateway IPs (the next-hop address servers use as their route to the internet).

Each attachment carries a set of gateway IPs. If you do not supply them they are auto-generated; when you do supply them they should belong to the same subnet as the LAN. Attaching a LAN by itself does not translate any traffic, you still need SNAT rules (` + "`" + `natgateway rule` + "`" + `) whose source subnet covers the servers on that LAN.`,
			TraverseChildren: true,
		},
	}

	natgatewayLanCmd.AddCommand(NatgatewayLanListCmd())
	natgatewayLanCmd.AddCommand(NatgatewayLanAddCmd())
	natgatewayLanCmd.AddCommand(NatgatewayLanRemoveCmd())

	return core.WithConfigOverride(natgatewayLanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
