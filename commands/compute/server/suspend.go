package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerSuspendCmd() *core.Command {
	suspend := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "suspend",
		ShortDesc: "Suspend (pause) a CUBE Server, deallocating compute while keeping its DAS data",
		LongDesc: `Use this command to suspend a CUBE Server. This operation is CUBE-only — ENTERPRISE, VCPU and GPU servers cannot be suspended (use ` + "`server stop`" + ` for those).

Suspending pauses the machine and deallocates its compute resources, but the virtual machine is NOT deleted: its NVMe Direct Attached Storage (DAS) boot volume, and the data on it, are retained. Resume it later with ` + "`server resume`" + ` to bring it back to a running state. Only servers whose type is CUBE are offered for completion here.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    "ionosctl compute server suspend --datacenter-id DATACENTER_ID -i SERVER_ID",
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerSuspend,
		InitClient: true,
	})
	suspend.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(viper.GetString(core.GetFlagName(suspend.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string][]string{
					"type": {"CUBE"},
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})

	return suspend
}
