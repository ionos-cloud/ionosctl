package pcc

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func PccGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Cross-Connect",
		LongDesc:   "Use this command to retrieve details about a specific Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id",
		Example:    `ionosctl pcc get --pcc-id PCC_ID`,
		PreCmdRun:  cloudapiv6cmds.PreRunPccId,
		CmdRun:     cloudapiv6cmds.RunPccGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
