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
			Use:              "targetgroup",
			Aliases:          []string{"tg"},
			Short:            "Target Group Operations",
			Long:             "The sub-commands of `ionosctl compute targetgroup` allow you to see information, to create, update, delete Target Groups.",
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
