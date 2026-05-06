package nic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allNicCols = []table.Column{
	{Name: "NicId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Dhcp", JSONPath: "properties.dhcp", Default: true},
	{Name: "LanId", JSONPath: "properties.lan", Default: true},
	{Name: "Ips", JSONPath: "properties.ips", Default: true},
	{Name: "IPv6Ips", JSONPath: "properties.ipv6Ips", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "FirewallActive", JSONPath: "properties.firewallActive"},
	{Name: "FirewallType", JSONPath: "properties.firewallType"},
	{Name: "DeviceNumber", JSONPath: "properties.deviceNumber"},
	{Name: "PciSlot", JSONPath: "properties.pciSlot"},
	{Name: "Mac", JSONPath: "properties.mac"},
	{Name: "DHCPv6", JSONPath: "properties.dhcpv6"},
	{Name: "IPv6CidrBlock", JSONPath: "properties.ipv6CidrBlock"},
}

func NicCmd() *core.Command {
	nicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Network Interfaces Operations",
			Long:             "The sub-commands of `ionosctl compute nic` allow you to create, list, get, update, delete NICs. To attach a NIC to a Load Balancer, use the Load Balancer command `ionosctl compute loadbalancer nic attach`.",
			TraverseChildren: true,
		},
	}
	nicCmd.AddColsFlag(allNicCols)

	nicCmd.AddCommand(NicListCmd())
	nicCmd.AddCommand(NicGetCmd())
	nicCmd.AddCommand(NicCreateCmd())
	nicCmd.AddCommand(NicUpdateCmd())
	nicCmd.AddCommand(NicDeleteCmd())

	return core.WithConfigOverride(nicCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
