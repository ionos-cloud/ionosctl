package gateway

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "ID", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "PublicKey", JSONPath: "metadata.publicKey", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "GatewayIP", JSONPath: "properties.gatewayIP", Default: true},
	{Name: "InterfaceIPv4", JSONPath: "properties.interfaceIPv4CIDR", Default: true},
	{Name: "InterfaceIPv6", JSONPath: "properties.interfaceIPv6CIDR", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId", Default: true},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId"},
	{Name: "ConnectionIPv4", JSONPath: "properties.connections.0.ipv4CIDR"},
	{Name: "ConnectionIPv6", JSONPath: "properties.connections.0.ipv6CIDR"},
	{Name: "InterfaceIP", JSONPath: "properties.interfaceIPv4CIDR"},
	{Name: "ListenPort", JSONPath: "properties.listenPort", Default: true},
	{Name: "Status", JSONPath: "metadata.status", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:   "gateway",
			Short: "Manage WireGuard VPN Gateways",
			Long: `Manage WireGuard VPN Gateways.

A gateway is the IONOS side of a WireGuard VPN. It attaches to one LAN in a datacenter and holds:
  - a key pair (you supply the PRIVATE key; the matching public key is shown so remote peers can trust it),
  - a public --gateway-ip (from an IPBlock) that peers dial, plus a private --connection-ip on the LAN,
  - a tunnel --interface-ip and a UDP --port it listens on (default 51820).

After the gateway is AVAILABLE, add remote devices with 'vpn wireguard peer create'.`,
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
