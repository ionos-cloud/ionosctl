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
		Namespace:  "ipconsumer",
		Resource:   "ipconsumer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IpConsumers",
		LongDesc:   "Use this command to get a list of Resources where each IP address from an IpBlock is being used.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    `ionosctl compute ipconsumer list --ipblock-id IPBLOCK_ID`,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpConsumersList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
