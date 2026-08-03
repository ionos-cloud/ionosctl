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

The ALB is the root of a four-level object tree:

  applicationloadbalancer   the balancer itself: which LAN it listens on (--listener-lan), which public/private IPs clients connect to (--ips), and which private LAN the backends live on (--target-lan)
    └── rule                a forwarding rule = a listener socket (protocol + IP + port), e.g. HTTP on 10.0.0.5:80
          └── httprule      HTTP rules evaluated in order within a listener; each matches conditions and then FORWARDs to a target group, returns a STATIC response, or issues a REDIRECT
    └── flowlog             captures ALB traffic metadata to an IONOS Object Storage (S3) bucket for auditing/troubleshooting

Backend servers are not registered on the ALB directly; they are grouped into a Target Group (a separate resource, see ` + "`" + `ionosctl compute target-group` + "`" + `) which a FORWARD httprule then points to. The full request path is therefore: client -> ALB listener IP (--ips) -> forwarding rule -> matching httprule -> target group -> backend server on the target LAN.

The sub-commands below let you create, list, get, update and delete Application Load Balancers.`,
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
