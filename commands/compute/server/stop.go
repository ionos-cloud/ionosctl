package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
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
		ShortDesc: "Stop a Server",
		LongDesc: `Use this command to stop a Server from a Virtual Data Center. The machine will be forcefully powered off, billing will cease, and the public IP, if one is allocated, will be deallocated.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

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
	stop.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server stop to be executed")
	stop.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server stop [seconds]")

	return stop
}
