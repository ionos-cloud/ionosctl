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
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Network Load Balancer Forwarding Rule Target Operations",
			Long:             "The sub-commands of `ionosctl compute networkloadbalancer rule target` allow you to add, list, update, remove Network Load Balancer Forwarding Rule Targets.",
			TraverseChildren: true,
		},
	}
	nlbRuleTargetCmd.AddColsFlag(allCols)

	nlbRuleTargetCmd.AddCommand(NlbRuleTargetListCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetAddCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetRemoveCmd())

	return core.WithConfigOverride(nlbRuleTargetCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
