package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/rule/httprule"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allAlbForwardingRuleCols = []table.Column{
	{Name: "ForwardingRuleId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Protocol", JSONPath: "properties.protocol", Default: true},
	{Name: "ListenerIp", JSONPath: "properties.listenerIp", Default: true},
	{Name: "ListenerPort", JSONPath: "properties.listenerPort", Default: true},
	{Name: "ClientTimeout", JSONPath: "properties.clientTimeout"},
	{Name: "ServerCertificates", JSONPath: "properties.serverCertificates", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ApplicationLoadBalancerRuleCmd() *core.Command {
	albRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "rule",
			Aliases: []string{"r", "forwardingrule"},
			Short:   "Application Load Balancer Forwarding Rule Operations",
			Long: `A forwarding rule is a listener on an Application Load Balancer: it binds a protocol, an IP and a port (e.g. HTTP on 192.0.2.10:80) that the balancer accepts client connections on. Each rule owns one listener socket; to serve both port 80 and 443 you create two forwarding rules.

The IP a rule listens on (--listener-ip) must be one of the ALB's own --ips on the listener LAN. For HTTPS listeners you attach one or more server certificates (--server-certificates, referencing certificates managed by the IONOS Certificate Manager) which the balancer presents during the TLS handshake.

A forwarding rule by itself only accepts connections; the actual request routing is defined by the HTTP rules inside it (` + "`" + `ionosctl compute alb rule httprule` + "`" + `), which match request attributes and then forward to a target group, return a static response, or redirect.

The sub-commands below let you create, list, get, update and delete forwarding rules.`,
			TraverseChildren: true,
		},
	}

	albRuleCmd.AddCommand(ApplicationLoadBalancerForwardingRuleListCmd())
	albRuleCmd.AddCommand(ApplicationLoadBalancerForwardingRuleGetCmd())
	albRuleCmd.AddCommand(ApplicationLoadBalancerForwardingRuleCreateCmd())
	albRuleCmd.AddCommand(ApplicationLoadBalancerForwardingRuleUpdateCmd())
	albRuleCmd.AddCommand(ApplicationLoadBalancerForwardingRuleDeleteCmd())

	albRuleCmd.AddCommand(httprule.AlbRuleHttpRuleCmd())

	return core.WithConfigOverride(albRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
