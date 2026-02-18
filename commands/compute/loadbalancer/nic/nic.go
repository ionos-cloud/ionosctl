package nic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultNicCols = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "IPv6Ips", "State"}
	allNicCols     = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "IPv6Ips", "State", "FirewallActive",
		"FirewallType", "DeviceNumber", "PciSlot", "Mac", "DHCPv6", "IPv6CidrBlock"}
)

func LoadBalancerNicCmd() *core.Command {
	loadbalancerNicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Load Balancer Nic Operations",
			Long:             "The sub-commands of `ionosctl loadbalancer nic` allow you to manage NICs on Load Balancers.",
			TraverseChildren: true,
		},
	}

	loadbalancerNicCmd.AddCommand(LoadBalancerNicAttachCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicListCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicGetCmd())
	loadbalancerNicCmd.AddCommand(LoadBalancerNicDetachCmd())

	return core.WithConfigOverride(loadbalancerNicCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
