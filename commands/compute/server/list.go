package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ServerListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "server",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Servers in a Virtual Data Center",
		LongDesc:   "Use this command to list the Servers in a Virtual Data Center (--datacenter-id), or across every datacenter in the contract with --all. The output shows each server's Type, sizing (Cores/RAM), CPU family and lifecycle state (VmState = running/deallocated; State = provisioning state).\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ServersFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id (unless --all is used)",
		Example:    "ionosctl compute server list --datacenter-id DATACENTER_ID\nionosctl compute server list --all\nionosctl compute server list --datacenter-id DATACENTER_ID --filters type=CUBE",
		PreCmdRun:  PreRunServerList,
		CmdRun:     RunServerList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	return cmd
}
