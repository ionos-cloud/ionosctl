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
			Use:              "gateway",
			Short:            "Manage Wireguard VPN Gateways",
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
