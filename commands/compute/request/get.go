package request

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func RequestGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "request",
		Resource:  "request",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get one request's status, method and target resources",
		LongDesc: `Use this command to inspect a single asynchronous request: its HTTP method, the resources it targets, and its current status (QUEUED / RUNNING / DONE / FAILED). When a request FAILED, the Message column explains why.

Required values to run command:

* Request Id (printed by the command that started the action)`,
		Example:    `ionosctl compute request get --request-id REQUEST_ID`,
		PreCmdRun:  PreRunRequestId,
		CmdRun:     RunRequestGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgRequestId, cloudapiv6.ArgIdShort, "", "The ID of the request to inspect", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
