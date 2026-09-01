package ipblock

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		ShortDesc: "Rename an existing IpBlock",
		LongDesc: `Update an existing IpBlock. Only the ` + "`" + `--name` + "`" + ` (friendly label) can be changed; the reserved addresses, their ` + "`" + `--location` + "`" + ` and the block ` + "`" + `--size` + "`" + ` are immutable. To change how many IPs you hold, reserve a new block (` + "`" + `ipblock create` + "`" + `) and delete the old one.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the IpBlock is back in AVAILABLE state.

Required values to run command:

* IpBlock Id`,
		Example:    "ionosctl compute ipblock update --ipblock-id IPBLOCK_ID --name new-label",
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New friendly label for the block. This is the only mutable property; it does not affect the reserved IP addresses")

	return cmd
}
