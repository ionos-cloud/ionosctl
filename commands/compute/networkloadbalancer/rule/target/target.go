package target

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "TargetIp", JSONPath: "ip", Default: true},
	{Name: "TargetPort", JSONPath: "port", Default: true},
	{Name: "Weight", JSONPath: "weight", Default: true},
	{Name: "Check", JSONPath: "healthCheck.check", Default: true},
	{Name: "CheckInterval", JSONPath: "healthCheck.checkInterval", Default: true},
	{Name: "Maintenance", JSONPath: "healthCheck.maintenance", Default: true},
}

func NlbRuleTargetCmd() *core.Command {
	nlbRuleTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "target",
			Aliases: []string{"t"},
			Short:   "Network Load Balancer Forwarding Rule Target Operations",
			Long: `A target is a backend VM that a forwarding rule distributes connections to, identified by its IP and port (--ip + --port) on the NLB's target LAN. A rule can have many targets; together they form the pool that the rule's balancing algorithm chooses from.

--weight controls each target's share of traffic relative to the other targets (higher weight = more connections); --check enables per-target TCP health probes so unhealthy VMs are pulled out of rotation; --maintenance lets you drain a target manually without deleting it.

Targets are addressed by IP+port, not by an ID, so add/remove use those values directly.`,
			TraverseChildren: true,
		},
	}
	nlbRuleTargetCmd.AddColsFlag(allCols)

	nlbRuleTargetCmd.AddCommand(NlbRuleTargetListCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetAddCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetRemoveCmd())

	return core.WithConfigOverride(nlbRuleTargetCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
