package targetgroup

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allTargetGroupCols = []table.Column{
	{Name: "TargetGroupId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Algorithm", JSONPath: "properties.algorithm", Default: true},
	{Name: "Protocol", JSONPath: "properties.protocol", Default: true},
	{Name: "CheckTimeout", JSONPath: "properties.healthCheck.timeout", Default: true},
	{Name: "CheckInterval", JSONPath: "properties.healthCheck.interval", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Retries", JSONPath: "properties.healthCheck.retries"},
	{Name: "Path", JSONPath: "properties.httpHealthCheck.path"},
	{Name: "Method", JSONPath: "properties.httpHealthCheck.method"},
	{Name: "MatchType", JSONPath: "properties.httpHealthCheck.matchType"},
	{Name: "Response", JSONPath: "properties.httpHealthCheck.response"},
	{Name: "Regex", JSONPath: "properties.httpHealthCheck.regex"},
	{Name: "Negate", JSONPath: "properties.httpHealthCheck.negate"},
}

func TargetGroupCmd() *core.Command {
	targetGroupCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "targetgroup",
			Aliases: []string{"tg"},
			Short:   "Target Group Operations",
			Long: `Manage Target Groups: reusable backend server pools referenced by Application Load Balancer (ALB) forwarding rules.

A Target Group is a location-independent, standalone resource in your contract (it is NOT nested under a datacenter or an ALB). It decouples the definition of "which backend servers exist and how they are health-checked" from "which load balancer routes traffic to them". An ALB forwarding rule with a FORWARD action (see ` + "`" + `ionosctl alb rule httprule` + "`" + `) references a Target Group by its ID; the same Target Group can be reused by multiple rules and multiple ALBs.

A Target Group is made of three parts:
  1. Distribution - the balancing --algorithm (how traffic is spread across targets) and --protocol (HTTP only).
  2. Health checking - a connection-level health check (--check-timeout, --check-interval, --retries) plus an application-level HTTP health check (--path, --method, --match-type, --response, --regex, --negate). A target is only sent traffic while it is considered healthy.
  3. Targets - the actual backend servers, each an IP + port + weight, managed with the ` + "`" + `target` + "`" + ` sub-commands (add/list/remove).

Balancing algorithms:
  - ROUND_ROBIN      Targets are served alternately, honoring their weights.
  - LEAST_CONNECTION The target with the fewest active connections is served next.
  - RANDOM           Targets are chosen by a consistent pseudo-random function.
  - SOURCE_IP        The same client IP is always routed to the same target (session stickiness by source address).

Note: --protocol only accepts HTTP.`,
			TraverseChildren: true,
		},
	}
	targetGroupCmd.AddColsFlag(allTargetGroupCols)

	targetGroupCmd.AddCommand(TargetGroupListCmd())
	targetGroupCmd.AddCommand(TargetGroupGetCmd())
	targetGroupCmd.AddCommand(TargetGroupCreateCmd())
	targetGroupCmd.AddCommand(TargetGroupUpdateCmd())
	targetGroupCmd.AddCommand(TargetGroupDeleteCmd())
	targetGroupCmd.AddCommand(TargetGroupTargetCmd())

	return core.WithConfigOverride(targetGroupCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
