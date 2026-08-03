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
			Use:     "request",
			Aliases: []string{"req"},
			Short:   "Inspect and wait on the async provisioning requests",
			Long: `The ` + "`ionosctl compute request`" + ` command lets you track the asynchronous requests that the IONOS CLOUD API creates for you.

Every provisioning action (create/update/delete a datacenter, server, volume, ...) is asynchronous: the API immediately returns a Request that then runs in the background. A request moves through these states:

* QUEUED  - accepted, not started yet
* RUNNING - being executed
* DONE    - completed successfully
* FAILED  - completed with an error (see the Message column for the reason)

The request ID appears in the output of the command that started the action (and in the ` + "`Location`" + ` response header). The ` + "`--wait-for-request`" + ` flag on create/update/delete commands is built on this same mechanism: it polls the request until it reaches DONE.

Use ` + "`list`" + ` to see recent requests, ` + "`get`" + ` to inspect one (its method, target resources and status), and ` + "`wait`" + ` to block until a specific request finishes. These commands are read-only; they do not change resources.`,
			TraverseChildren: true,
		},
	}
	reqCmd.AddColsFlag(allRequestCols)

	reqCmd.AddCommand(RequestListCmd())
	reqCmd.AddCommand(RequestGetCmd())
	reqCmd.AddCommand(RequestWaitCmd())

	return core.WithConfigOverride(reqCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
