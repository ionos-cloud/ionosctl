package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/rule/target"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "ForwardingRuleId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Algorithm", JSONPath: "properties.algorithm", Default: true},
	{Name: "Protocol", JSONPath: "properties.protocol", Default: true},
	{Name: "ListenerIp", JSONPath: "properties.listenerIp", Default: true},
	{Name: "ListenerPort", JSONPath: "properties.listenerPort", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "ClientTimeout", JSONPath: "properties.healthCheck.clientTimeout"},
	{Name: "ConnectTimeout", JSONPath: "properties.healthCheck.connectTimeout"},
	{Name: "TargetTimeout", JSONPath: "properties.healthCheck.targetTimeout"},
	{Name: "Retries", JSONPath: "properties.healthCheck.retries"},
}

func NetworkloadbalancerRuleCmd() *core.Command {
	nlbRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Network Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl compute nlb rule` allow you to create, list, get, update, delete Network Load Balancer Forwarding Rules.",
			TraverseChildren: true,
		},
	}
	nlbRuleCmd.AddColsFlag(allCols)

	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleListCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleGetCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleCreateCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleUpdateCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleDeleteCmd())

	nlbRuleCmd.AddCommand(target.NlbRuleTargetCmd())

	return core.WithConfigOverride(nlbRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
