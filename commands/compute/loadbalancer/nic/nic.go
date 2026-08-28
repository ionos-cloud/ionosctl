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
			Use:     "nic",
			Aliases: []string{"n"},
			Short:   "Load Balancer NIC Operations",
			Long: `The sub-commands of ` + "`" + `ionosctl compute loadbalancer nic` + "`" + ` manage which server NICs are attached to a (basic) Load Balancer.

A Load Balancer distributes incoming traffic across a pool of backend servers. That pool is defined by NICs: each NIC you attach registers its server as a balancer backend. Once attached, the NIC inherits the Load Balancer's public IPv4 address, so all backends share the same balanced, public-facing IP; traffic sent to that IP is spread across the attached NICs.

Relationships:
  * A NIC belongs to exactly one server, which belongs to one Data Center. Identify a NIC by --datacenter-id + --nic-id (optionally --server-id, used only to speed up completion).
  * A NIC may be attached to only one Load Balancer at a time.
  * Detaching a NIC removes it from the balancer pool; the NIC and its server are left intact.

Note: this is the original/basic IONOS Load Balancer, which balances by server NIC. It is distinct from the newer, feature-rich 'applicationloadbalancer' (Layer 7 / HTTP) and 'networkloadbalancer' (Layer 4) resources; for new deployments prefer those.`,
			TraverseChildren: true,
		},
	}

	loadbalancerNicCmd.AddCommand(LoadBalancerNicAttachCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicListCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicGetCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicDetachCmd())

	return core.WithConfigOverride(loadbalancerNicCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
