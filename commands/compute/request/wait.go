package request

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func RequestWaitCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "request",
		Resource:  "request",
		Verb:      "wait",
		Aliases:   []string{"w"},
		ShortDesc: "Block until a request reaches DONE (or FAILED)",
		LongDesc: `Use this command to block the terminal until a specific asynchronous request finishes. The CLI polls the request's status endpoint and returns once the request reaches a terminal state: it prints the request on success (DONE) or returns an error on FAILED.

This is the same polling that ` + "`--wait-for-request`" + ` performs inline on create/update/delete commands; use ` + "`request wait`" + ` when you already have a request ID and want to wait after the fact (e.g. in a script).

Use the global ` + "`--timeout`" + ` option to cap how long to wait, in seconds (default 600). If the request has not finished within the timeout, the command aborts with a timeout error.

Required values to run command:

* Request Id`,
		Example: `# Wait up to the default 600s for a request to finish
ionosctl compute request wait --request-id REQUEST_ID

# Give up after 120 seconds
ionosctl compute request wait --request-id REQUEST_ID --timeout 120`,
		PreCmdRun:  PreRunRequestId,
		CmdRun:     RunRequestWait,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgRequestId, cloudapiv6.ArgIdShort, "", "The ID of the request to wait on", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
