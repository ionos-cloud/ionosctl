package label

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LabelRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "label",
		Resource:   "label",
		Verb:       "remove",
		Aliases:    []string{"delete", "del", "r", "rm"},
		ShortDesc:  "Remove a Label from a Resource",
		LongDesc:   "Use this command to remove a Label from a Resource.\n\nRequired values to run command:\n\n* Resource Type\n* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id, Image ID, or Snapshot Id\n* Label Key",
		Example:    "ionosctl label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY",
		PreCmdRun:  cloudapiv6cmds.PreRunResourceTypeLabelKeyRemove,
		CmdRun:     cloudapiv6cmds.RunLabelRemove,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", cloudapiv6.LabelKey, core.RequiredFlagOption())
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId)
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId)
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId)
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId)
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", cloudapiv6.SnapshotId)
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", cloudapiv6.ImageId+"(note: only private images supported)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(
			func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
				return request.Filter("public", "false")
			}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Type of resource to remove labels from", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Labels")

	return cmd
}
