package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		LongDesc: `Use this command to update the name and/or description of a specific Cross-Connect. This only changes the Cross-Connect's own metadata; it does not attach or detach LANs. To change which LANs are peered, use ` + "`" + `ionosctl compute lan update --pcc-id` + "`" + ` (attach) or clear a LAN's Cross-Connect (detach).

Required values to run command:

* Pcc Id`,
		Example:    `ionosctl compute pcc update --pcc-id PCC_ID --description "New description"`,
		PreCmdRun:  PreRunPccId,
		CmdRun:     RunPccUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The new name for the Cross-Connect. Leave unset to keep the current name")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "The new description for the Cross-Connect. Leave unset to keep the current description")
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", "The unique ID of the Cross-Connect to update", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
