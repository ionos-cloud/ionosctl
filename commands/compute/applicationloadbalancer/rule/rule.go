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
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Application Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl compute alb rule` allow you to create, list, get, update, delete Application Load Balancer Forwarding Rules.",
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
