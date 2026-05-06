package applicationloadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allApplicationLoadBalancerCols = []table.Column{
	{Name: "ApplicationLoadBalancerId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "ListenerLan", JSONPath: "properties.listenerLan", Default: true},
	{Name: "Ips", JSONPath: "properties.ips", Default: true},
	{Name: "TargetLan", JSONPath: "properties.targetLan", Default: true},
	{Name: "PrivateIps", JSONPath: "properties.lbPrivateIps", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "DatacenterId", JSONPath: "href"},
}

func ApplicationLoadBalancerCmd() *core.Command {
	applicationloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "applicationloadbalancer",
			Aliases:          []string{"alb"},
			Short:            "Application Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl compute applicationloadbalancer` allow you to create, list, get, update, delete Application Load Balancers.",
			TraverseChildren: true,
		},
	}
	applicationloadbalancerCmd.AddColsFlag(allApplicationLoadBalancerCols)

	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerListCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerGetCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerCreateCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerUpdateCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerDeleteCmd())

	applicationloadbalancerCmd.AddCommand(rule.ApplicationLoadBalancerRuleCmd())
	applicationloadbalancerCmd.AddCommand(flowlog.ApplicationLoadBalancerFlowLogCmd())

	return core.WithConfigOverride(applicationloadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
