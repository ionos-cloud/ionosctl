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

func LoadBalancerNicCmd() *core.Command {
	loadbalancerNicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Load Balancer Nic Operations",
			Long:             "The sub-commands of `ionosctl compute loadbalancer nic` allow you to manage NICs on Load Balancers.",
			TraverseChildren: true,
		},
	}

	loadbalancerNicCmd.AddCommand(LoadBalancerNicAttachCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicListCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicGetCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicDetachCmd())

	return core.WithConfigOverride(loadbalancerNicCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
