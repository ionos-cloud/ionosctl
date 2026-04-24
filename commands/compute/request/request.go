package request

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allRequestCols = []table.Column{
	{Name: "RequestId", JSONPath: "id", Default: true},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate", Default: true},
	{Name: "CreatedBy", JSONPath: "metadata.createdBy"},
	{Name: "Method", JSONPath: "properties.method", Default: true},
	{Name: "Status", JSONPath: "metadata.requestStatus.metadata.status", Default: true},
	{Name: "Message", JSONPath: "metadata.requestStatus.metadata.message", Default: true},
	{Name: "Url", JSONPath: "properties.url"},
	{Name: "Body", JSONPath: "properties.body"},
	{Name: "Targets", Default: true, Format: func(item map[string]any) any {
		targets := table.Navigate(item, "metadata.requestStatus.metadata.targets")
		if targets == nil {
			return ""
		}
		targetsSlice, ok := targets.([]any)
		if !ok {
			return ""
		}
		var parts []string
		for _, t := range targetsSlice {
			tMap, ok := t.(map[string]any)
			if !ok {
				continue
			}
			targetMap, ok := tMap["target"].(map[string]any)
			if !ok {
				continue
			}
			id, _ := targetMap["id"].(string)
			typ, _ := targetMap["type"].(string)
			if id != "" && typ != "" {
				parts = append(parts, fmt.Sprintf("%s (%s)", id, typ))
			}
		}
		return strings.Join(parts, ", ")
	}},
}

func RequestCmd() *core.Command {
	reqCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "request",
			Aliases:          []string{"req"},
			Short:            "Request Operations",
			Long:             "The sub-commands of `ionosctl compute request` allow you to see information about requests on your account. With the `ionosctl compute request` command, you can list, get or wait for a Request.",
			TraverseChildren: true,
		},
	}
	reqCmd.AddColsFlag(allRequestCols)

	reqCmd.AddCommand(RequestListCmd())
	reqCmd.AddCommand(RequestGetCmd())
	reqCmd.AddCommand(RequestWaitCmd())

	return core.WithConfigOverride(reqCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
