package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func PccDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "pcc",
		Resource:  "pcc",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Cross-Connect",
		LongDesc: `Use this command to delete a Cross-Connect. A Cross-Connect can only be deleted once no LANs are peered through it; detach every LAN first (clear each LAN's Cross-Connect) or the request will be rejected.

Required values to run command:

* Pcc Id`,
		Example: `# Delete a single Cross-Connect and wait for completion
ionosctl compute pcc delete --pcc-id PCC_ID --wait

# Delete every Cross-Connect on the contract
ionosctl compute pcc delete --all`,
		PreCmdRun:  PreRunPccDelete,
		CmdRun:     RunPccDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, cloudapiv6.ArgIdShort, "", "The unique ID of the Cross-Connect to delete", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Cross-Connects on the contract. Mutually exclusive with --pcc-id")

	return cmd
}
