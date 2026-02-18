package template

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TemplateGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "template",
		Resource:   "template",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Template",
		LongDesc:   "Use this command to get information about a specified Template.\n\nRequired values to run command:\n\n* Template Id",
		Example:    `ionosctl template get -i TEMPLATE_ID`,
		PreCmdRun:  cloudapiv6cmds.PreRunTemplateId,
		CmdRun:     cloudapiv6cmds.RunTemplateGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTemplateId, cloudapiv6.ArgIdShort, "", cloudapiv6.TemplateId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
