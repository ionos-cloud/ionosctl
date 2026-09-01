package ipblock

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func IpBlockDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete (release) an IpBlock",
		LongDesc: `Release a reserved IpBlock, returning all of its IPs to the pool and stopping billing for them.

An IP that is still assigned to a consumer (a NIC, NAT gateway, load balancer or IP-failover group) cannot be released - detach the IP from that resource first, then delete the block. Use ` + "`" + `ionosctl compute ipconsumer list --ipblock-id <id>` + "`" + ` to see what is still holding an IP.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the deletion completes. Use ` + "`" + `--force` + "`" + ` to skip the confirmation prompt.

Required values to run command:

* IpBlock Id`,
		Example:    "ionosctl compute ipblock delete --ipblock-id IPBLOCK_ID --wait",
		PreCmdRun:  PreRunIpBlockDelete,
		CmdRun:     RunIpBlockDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Release every IpBlock on the contract (only those whose IPs are not in use). Use instead of --ipblock-id")

	return cmd
}
