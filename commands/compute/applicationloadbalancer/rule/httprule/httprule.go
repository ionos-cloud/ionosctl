package httprule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allAlbRuleHttpRuleCols = []table.Column{
	{Name: "Name", JSONPath: "name", Default: true},
	{Name: "Type", JSONPath: "type", Default: true},
	{Name: "TargetGroupId", JSONPath: "targetGroup", Default: true},
	{Name: "DropQuery", JSONPath: "dropQuery", Default: true},
	{Name: "Location", JSONPath: "location"},
	{Name: "StatusCode", JSONPath: "statusCode"},
	{Name: "ResponseMessage", JSONPath: "responseMessage"},
	{Name: "ContentType", JSONPath: "contentType"},
	{Name: "Condition", JSONPath: "conditions", Default: true},
}

func AlbRuleHttpRuleCmd() *core.Command {
	albRuleHttpRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "httprule",
			Aliases:          []string{"http"},
			Short:            "Application Load Balancer Forwarding Rule Http Rule Operations",
			Long:             "The sub-commands of `ionosctl compute alb rule httprule` allow you to add, list, update, remove Application Load Balancer Forwarding Rule Http Rules.",
			TraverseChildren: true,
		},
	}
	albRuleHttpRuleCmd.AddColsFlag(allAlbRuleHttpRuleCols)

	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleListCmd())
	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleAddCmd())
	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleRemoveCmd())

	return core.WithConfigOverride(albRuleHttpRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
