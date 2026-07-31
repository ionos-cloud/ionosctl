package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerStopCmd() *core.Command {
	stop := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "stop",
		Aliases:   []string{"off"},
		ShortDesc: "Stop (deallocate) a Server; halts compute billing and releases its public IP",
		LongDesc: `Use this command to stop a Server in a Virtual Data Center. The machine is forcefully powered off (this is NOT a graceful OS shutdown) and its compute resources (cores + RAM) are DEALLOCATED, so per-server compute billing ceases. Attached storage Volumes are kept and continue to be billed; the data on them is preserved.

Any public IP allocated to the Server is released on stop, and a new one is typically assigned when you ` + "`server start`" + ` it again — so do not rely on the IP staying the same across a stop/start cycle.

Note the difference from ` + "`server suspend`" + `: stop applies to any server type and fully deallocates compute, whereas suspend is a CUBE-only pause. Use ` + "`server start`" + ` to power a stopped server back on.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    "ionosctl compute server stop --datacenter-id DATACENTER_ID --server-id SERVER_ID",
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerStop,
		InitClient: true,
	})
	stop.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(stop.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return stop
}
