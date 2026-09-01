package label

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LabelRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "label",
		Resource:  "label",
		Verb:      "remove",
		Aliases:   []string{"delete", "del", "r", "rm"},
		ShortDesc: "Remove a Label from a Resource",
		LongDesc: `Use this command to remove a Label from a Resource. Select the resource with --resource-type plus its id flag(s) (see the same pairing as ` + "`label add`" + `) and identify the label by --label-key.

Use --all to remove every label. With --resource-type and its id flag(s) it removes all labels on that one resource; with no --resource-type it iterates over labels of ALL resources under your account (prompting for each unless --force is given).

Required values to run command:

* Resource Type
* Resource Id(s) for that type
* Label Key`,
		Example: `# Remove one label from a datacenter
ionosctl compute label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key env

# Remove all labels from a server
ionosctl compute label remove --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID --all`,
		PreCmdRun:  PreRunResourceTypeLabelKeyRemove,
		CmdRun:     RunLabelRemove,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "The key of the label to remove from the resource", core.RequiredFlagOption())
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The Data Center Id. Required for --resource-type datacenter, server and volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", "The Server Id (also needs --datacenter-id). Required for --resource-type server")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", "The Volume Id (also needs --datacenter-id). Required for --resource-type volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", "The IpBlock Id. Required for --resource-type ipblock")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", "The Snapshot Id. Required for --resource-type snapshot")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "The Image Id (private images only). Required for --resource-type image")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(
			func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
				return request.Filter("public", "false")
			}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "The kind of resource to remove the label from. Determines which id flag(s) are required", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all labels: on the selected resource when --resource-type is given, otherwise across every labeled resource on the account")

	return cmd
}
