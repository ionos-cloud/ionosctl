package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/rule/httprule"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultAlbForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Protocol", "ListenerIp", "ListenerPort", "ServerCertificates", "State"}
	allAlbForwardingRuleCols     = []string{"ForwardingRuleId", "Name", "Protocol", "ListenerIp", "ListenerPort", "ClientTimeout", "ServerCertificates", "State"}
)

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
