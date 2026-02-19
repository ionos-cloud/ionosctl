package ipblock

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func IpBlockUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an IpBlock",
		LongDesc: `Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* IpBlock Id`,
		Example:    "ionosctl compute ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME",
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the IpBlock")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IpBlock update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock update [seconds]")

	return cmd
}
