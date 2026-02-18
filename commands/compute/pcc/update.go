package pcc

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func PccUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "pcc",
		Resource:  "pcc",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Cross-Connect",
		LongDesc: `Use this command to update details about a specific Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id`,
		Example:    `ionosctl pcc update --pcc-id PCC_ID --description DESCRIPTION`,
		PreCmdRun:  cloudapiv6cmds.PreRunPccId,
		CmdRun:     cloudapiv6cmds.RunPccUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Cross-Connect")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "The description for the Cross-Connect")
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Cross-Connect update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Cross-Connect update [seconds]")

	return cmd
}
