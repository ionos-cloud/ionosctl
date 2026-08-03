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
			Use:     "networkloadbalancer",
			Aliases: []string{"nlb"},
			Short:   "Network Load Balancer Operations",
			Long: `A Network Load Balancer (NLB) is a layer-4 (TCP) load balancer that lives inside a Virtual Data Center and spreads inbound connections across a pool of backend VMs.

Because it works at the transport layer it forwards raw TCP connections without inspecting their contents. It does not terminate TLS, read HTTP headers, or route by URL/host - if you need those layer-7 features, use the Application Load Balancer (ALB) instead. The NLB is the right choice for high-throughput, protocol-agnostic traffic (databases, custom TCP services, TLS pass-through, etc.).

The NLB sits between two LANs in the same data center:
  - the listener LAN (--listener-lan), where clients connect. Its addresses (--ips) are the public/customer-reserved IPs for a public NLB, or private IPs for a private NLB.
  - the target LAN (--target-lan), the private network where the balanced backend VMs live. The NLB reaches them over --lb-private-ips.

Resource hierarchy:
  Network Load Balancer -> forwarding rule (` + "`nlb rule`" + `) -> target (` + "`nlb rule target`" + `)

A forwarding rule binds a listener IP+port and picks a balancing algorithm; each target is a backend VM (IP+port) that the rule distributes connections to. Optionally attach a flowlog (` + "`nlb flowlog`" + `) to stream connection logs to an S3 bucket.`,
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
