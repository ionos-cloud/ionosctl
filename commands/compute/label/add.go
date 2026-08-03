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

func LabelAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "label",
		Resource:  "label",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Label to a Resource",
		LongDesc: `Use this command to add (or overwrite) a Label on a specific Resource.

Pick the target with --resource-type and supply the matching id flag(s):
  * datacenter -> --datacenter-id
  * server     -> --datacenter-id and --server-id
  * volume     -> --datacenter-id and --volume-id
  * snapshot   -> --snapshot-id
  * ipblock    -> --ipblock-id
  * image      -> --image-id (private images only)

Adding a key that already exists on the resource overwrites its value (there is one value per key per resource).

Required values to run command:

* Resource Type
* Resource Id(s) as listed above
* Label Key
* Label Value`,
		Example: `# Label a server
ionosctl compute label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID --label-key env --label-value prod

# Label a datacenter
ionosctl compute label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key team --label-value payments`,
		PreCmdRun:  PreRunResourceTypeLabelKeyValue,
		CmdRun:     RunLabelAdd,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "The label key. Unique per resource; adding an existing key overwrites its value", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", "The label value to store under the key", core.RequiredFlagOption())
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
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "The Image Id (private images only; public images cannot be labeled). Required for --resource-type image")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(
			func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
				return request.Filter("public", "false")
			}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgResourceType, "", "", allowedValues, "The kind of resource to label. Determines which id flag(s) are required (see command description)", core.RequiredFlagOption())

	return cmd
}
