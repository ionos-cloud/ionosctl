package cdrom

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
)

func ServerCdromAttachCmd() *core.Command {
	attachCdrom := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "cdrom",
		Verb:      "attach",
		Aliases:   []string{"a"},
		ShortDesc: "Attach a CD-ROM to a Server",
		LongDesc: `Use this command to attach a CD-ROM to an existing Server.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id`,
		Example:    "ionosctl compute server cdrom attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait",
		PreCmdRun:  PreRunDcServerCdromIds,
		CmdRun:     RunServerCdromAttach,
		InitClient: true,
	})
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgCdromId, cloudapiv6.ArgIdShort, "", cloudapiv6.CdromId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for CDROM images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				attachCdrom.Flags().String(cloudapiv6.ArgDataCenterId)).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "CDROM")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(attachCdrom.Flags().String(cloudapiv6.ArgDataCenterId)), cobra.ShellCompDirectiveNoFileComp
	})

	return attachCdrom
}
