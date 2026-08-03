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
			Use:     "rule",
			Aliases: []string{"r", "forwardingrule"},
			Short:   "Network Load Balancer Forwarding Rule Operations",
			Long: `A forwarding rule tells a Network Load Balancer how to accept and distribute traffic. Each rule listens on one IP+port pair on the NLB's listener LAN (--listener-ip + --listener-port, TCP) and forwards accepted connections to the rule's targets - the backend VMs added with ` + "`nlb rule target add`" + `.

The rule's --algorithm decides which target each new connection goes to:
  ROUND_ROBIN      cycle through targets in order (default).
  LEAST_CONNECTION send to the target with the fewest active connections.
  RANDOM           pick a target at random.
  SOURCE_IP        hash the client IP so a given client sticks to one target (basic session affinity).

Each rule also carries health-check timeouts (client/connect/target) and a retry count that govern how the NLB tolerates slow or failing targets. A rule with no targets accepts connections but has nothing to forward them to.`,
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
