package template

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TemplateGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "template",
		Resource:  "template",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get one CUBE template by ID",
		LongDesc: `Use this command to get the fixed cores, RAM, storage size and GPUs of a single CUBE template.

Required values to run command:

* Template Id (list them with ` + "`ionosctl compute template list`" + `)`,
		Example:    `ionosctl compute template get -i TEMPLATE_ID`,
		PreCmdRun:  PreRunTemplateId,
		CmdRun:     RunTemplateGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTemplateId, cloudapiv6.ArgIdShort, "", "The ID of the CUBE template to inspect", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
