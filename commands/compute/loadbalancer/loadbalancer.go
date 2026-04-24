package loadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/loadbalancer/nic"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allLoadbalancerCols = []table.Column{
	{Name: "LoadBalancerId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Dhcp", JSONPath: "properties.dhcp", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Ip", JSONPath: "properties.ip"},
	{Name: "DatacenterId", JSONPath: "href"},
}

func LoadBalancerCmd() *core.Command {
	loadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "loadbalancer",
			Aliases:          []string{"lb"},
			Short:            "Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl compute loadbalancer` manage your Load Balancers on your account. With Load Balancers you can distribute traffic between your servers.",
			TraverseChildren: true,
		},
	}
	loadbalancerCmd.AddColsFlag(allLoadbalancerCols)

	loadbalancerCmd.AddCommand(LoadBalancerListCmd())
	loadbalancerCmd.AddCommand(LoadBalancerGetCmd())
	loadbalancerCmd.AddCommand(LoadBalancerCreateCmd())
	loadbalancerCmd.AddCommand(LoadBalancerUpdateCmd())
	loadbalancerCmd.AddCommand(LoadBalancerDeleteCmd())
	loadbalancerCmd.AddCommand(nic.LoadBalancerNicCmd())

	return core.WithConfigOverride(loadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
