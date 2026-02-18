package cdrom

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerCdromDetachCmd() *core.Command {
	detachCdrom := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "cdrom",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a CD-ROM from a Server",
		LongDesc: `This will detach the CD-ROM from the Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id`,
		Example:    "ionosctl server cdrom detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request --force",
		PreCmdRun:  PreRunDcServerCdromDetach,
		CmdRun:     RunServerCdromDetach,
		InitClient: true,
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgCdromId, cloudapiv6.ArgIdShort, "", cloudapiv6.CdromId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedCdromsIds(viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for CD-ROM detachment to be executed")
	detachCdrom.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for CD-ROM detachment [seconds]")
	detachCdrom.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach all CD-ROMS from a Server.")

	return detachCdrom
}
