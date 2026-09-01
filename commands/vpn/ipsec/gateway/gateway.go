package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "ID", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "GatewayIP", JSONPath: "properties.gatewayIP", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId", Default: true},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId"},
	{Name: "ConnectionIPv4", JSONPath: "properties.connections.0.ipv4CIDR"},
	{Name: "ConnectionIPv6", JSONPath: "properties.connections.0.ipv6CIDR"},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Status", JSONPath: "metadata.status", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "gateway",
			Short: "Manage IPSec VPN Gateways",
			Long: `Manage IPSec VPN Gateways.

A gateway is the IONOS side of an IPSec VPN. It attaches to one LAN in a datacenter, takes a public --gateway-ip (from an IPBlock) that remote sites connect to, and a private --connection-ip on the LAN. --version pins the IKE version (IKEv2).

The gateway itself carries no crypto or remote-host settings — those live on the tunnels. After it is AVAILABLE, add one tunnel per remote site with 'vpn ipsec tunnel create'.`,
			Aliases:          []string{"g", "gw"},
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Get())
	cmd.AddCommand(Update())

	return cmd
}
