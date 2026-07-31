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
			Use:     "nic",
			Aliases: []string{"n"},
			Short:   "Network Interfaces Operations",
			Long: `The sub-commands of ` + "`ionosctl compute nic`" + ` let you manage NICs (Network Interface Cards) on IONOS Cloud servers.

A NIC is the virtual network adapter that connects a server to a LAN inside a Virtual Data Center. The resource hierarchy is:

  Data Center (--datacenter-id)  ->  Server (--server-id)  ->  NIC (--nic-id)  ->  LAN (--lan-id)

Every NIC lives on exactly one server and is attached to exactly one LAN; that LAN determines the network the server can reach. A single server may have several NICs, each on a different LAN, so a server can straddle multiple networks (for example one public-facing LAN and one internal LAN).

Addressing: a NIC receives its IP either automatically from the LAN's DHCP server (--dhcp=true, the default) or from IPs you assign explicitly (--ips). Explicitly assigned public IPs must come from reserved IP blocks; leaving --ips empty lets IONOS assign an address automatically. Private-range addresses (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) can be assigned manually on private LANs.

Firewall: each NIC has its own firewall toggle (--firewall-active). When active, all incoming traffic is blocked except what firewall rules on the NIC explicitly allow; when inactive, traffic flows straight to the NIC and any rules are ignored. --firewall-type controls the direction of traffic those rules apply to (INGRESS/EGRESS/BIDIRECTIONAL).

To attach a NIC to a Load Balancer instead, use ` + "`ionosctl compute loadbalancer nic attach`" + `.`,
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
