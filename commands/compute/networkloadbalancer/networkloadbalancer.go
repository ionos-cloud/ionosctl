package networkloadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "NetworkLoadBalancerId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "ListenerLan", JSONPath: "properties.listenerLan", Default: true},
	{Name: "Ips", JSONPath: "properties.ips", Default: true},
	{Name: "TargetLan", JSONPath: "properties.targetLan", Default: true},
	{Name: "LbPrivateIps", JSONPath: "properties.lbPrivateIps", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "DatacenterId", JSONPath: "href"},
}

func NetworkloadbalancerCmd() *core.Command {
	networkloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "networkloadbalancer",
			Aliases:          []string{"nlb"},
			Short:            "Network Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl compute networkloadbalancer` allow you to create, list, get, update, delete Network Load Balancers.",
			TraverseChildren: true,
		},
	}
	networkloadbalancerCmd.AddColsFlag(allCols)

	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerListCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerGetCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerCreateCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerUpdateCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerDeleteCmd())

	networkloadbalancerCmd.AddCommand(flowlog.NetworkloadbalancerFlowLogCmd())
	networkloadbalancerCmd.AddCommand(rule.NetworkloadbalancerRuleCmd())

	return core.WithConfigOverride(networkloadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
