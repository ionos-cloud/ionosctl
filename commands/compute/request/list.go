package request

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func RequestListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "request",
		Resource:  "request",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List recent provisioning requests and their status",
		LongDesc:  "Use this command to list the asynchronous requests on your account, most useful for reviewing recent activity or finding a failed request. The Status column shows QUEUED / RUNNING / DONE / FAILED, and Targets shows the resources each request acted on.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.RequestsFiltersUsage(),
		Example: `# List all requests
ionosctl compute request list

# Show only DONE/FAILED status and target columns
ionosctl compute request list --cols RequestId,Method,Status,Message,Targets`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunRequestList,
		InitClient: true,
	})
	cmd.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show only the N most recent requests (sorted by creation date, newest first). If unset, all requests are printed", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgMethod, "", "", "Show only requests with this HTTP method. Accepts POST/PUT/PATCH/DELETE, or the aliases CREATE (=POST) and UPDATE (=PUT+PATCH)", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"POST", "PUT", "DELETE", "PATCH", "CREATE", "UPDATE"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
