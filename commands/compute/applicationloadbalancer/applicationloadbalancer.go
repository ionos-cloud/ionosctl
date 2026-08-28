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
			Use:     "applicationloadbalancer",
			Aliases: []string{"alb"},
			Short:   "Application Load Balancer Operations",
			Long: `An Application Load Balancer (ALB) is a layer-7 (HTTP/HTTPS) load balancer that lives inside a Virtual Data Center. It terminates client connections on one LAN and distributes requests across backend targets on another LAN, making application-aware routing decisions based on the contents of each HTTP request (path, headers, method, host, etc.). This is in contrast to the Network Load Balancer (NLB), which operates at layer 4 (TCP) and only forwards packets by IP/port.

Requests flow through a forwarding rule (a listener socket of protocol + IP + port), whose ordered HTTP rules match conditions and then FORWARD to a target group, return a STATIC response, or issue a REDIRECT. Backend servers are not registered on the ALB directly; they are grouped into a Target Group (a separate resource, see ` + "`" + `ionosctl compute target-group` + "`" + `) which a FORWARD http-rule then points to. The full request path is therefore: client -> ALB listener IP (--ips) -> forwarding rule -> matching http-rule -> target group -> backend server on the target LAN.`,
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
