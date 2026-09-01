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

var (
	allowedValues = []string{cloudapiv6.DatacenterResource, cloudapiv6.VolumeResource, cloudapiv6.ServerResource,
		cloudapiv6.SnapshotResource, cloudapiv6.IpBlockResource, cloudapiv6.ImageResource}
)

func LabelListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "label",
		Resource:  "label",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Labels from Resources",
		LongDesc: `Use this command to list Labels.

Without --resource-type it lists every label across ALL resources under your account (each row shows key, value, resource type, resource id and URN). To scope to one resource, pass --resource-type plus that resource's id flag(s), e.g. --resource-type server needs --datacenter-id and --server-id.

You can filter the results using ` + "`--filters`" + ` option. Use the following format to set filters: ` + "`--filters KEY1=VALUE1,KEY2=VALUE2`" + `.
` + completer.LabelsFiltersUsage(),
		Example: `# All labels on the account
ionosctl compute label list

# Labels on one datacenter
ionosctl compute label list --resource-type datacenter --datacenter-id DATACENTER_ID`,
		PreCmdRun:  PreRunLabelList,
		CmdRun:     RunLabelList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", "The Data Center Id. Required with --resource-type datacenter, server and volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", "The Server Id (also needs --datacenter-id). Used with --resource-type server")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", "The Volume Id (also needs --datacenter-id). Used with --resource-type volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", "The IpBlock Id. Used with --resource-type ipblock")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, "", "", "The Snapshot Id. Used with --resource-type snapshot")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "The Image Id (private images only). Used with --resource-type image")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(
			func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
				return request.Filter("public", "false")
			}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "Scope the listing to one resource kind (also pass that kind's id flag(s)). If not given, labels across all resources are listed", core.RequiredFlagOption())

	return cmd
}
