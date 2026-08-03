package ipconsumer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func IpconsumerListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipconsumer",
		Resource:  "ipconsumer",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List the resources consuming each IP in an IP block",
		LongDesc: `Use this command to list, for every IP address in a reserved IP block, the resource that is currently using it: the NIC (and its MAC), the owning server and datacenter, and any Kubernetes cluster / node pool.

An empty result means none of the block's addresses are in use, so the block can be safely released.

Required values to run command:

* IpBlock Id (get it from ` + "`ionosctl compute ipblock list`" + `)`,
		Example: `# List consumers of every IP in a block
ionosctl compute ipconsumer list --ipblock-id IPBLOCK_ID

# Show only the IP, server and datacenter columns
ionosctl compute ipconsumer list --ipblock-id IPBLOCK_ID --cols Ip,ServerName,DatacenterName`,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpConsumersList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", "The ID of the reserved IP block whose addresses you want to inspect", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
