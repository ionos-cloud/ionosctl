package volume

import (
	"context"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VolumeCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Volume",
		LongDesc: `Use this command to create a Volume on your account, within a Data Center. This will NOT attach the Volume to a Server. Please see the Servers commands for details on how to attach storage Volumes. You can specify the name, size, type, licence type, availability zone, image and other properties for the object.

NNote: The Licence Type has a default value, but if Image ID or Image Alias is supplied, then Licence Type will be automatically set. The Image Password or SSH Keys attributes can be defined when creating a Volume that uses an Image ID or Image Alias of an IONOS public Image. You may wish to set a valid value for Image Password even when using SSH Keys so that it is possible to authenticate with a password when using the remote console feature of the DCD.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example: `ionosctl compute volume create --datacenter-id DATACENTER_ID --name NAME

ionosctl compute volume create --datacenter-id DATACENTER_ID --name NAME --image-alias IMAGE_ALIAS --ssh-keys-path "SSH_KEY_PATH1,SSH_KEY_PATH2"`,
		PreCmdRun:  PreRunVolumeCreate,
		CmdRun:     RunVolumeCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Volume", "Name of the Volume")
	cmd.AddStringFlag(cloudapiv6.ArgSize, cloudapiv6.ArgSizeShort, strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Volume in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "The bus type of the Volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "LINUX", constants.EnumLicenceType, "Licence Type of the Volume")
	cmd.AddStringFlag(constants.FlagType, "", "HDD", "Type of the Volume")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, "", "", "The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "The Image Id or Snapshot Id to be used as template for the new Volume. A password or SSH Key need to be set")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		imageIds := completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "HDD")
		})

		snapshotIds := completer.SnapshotIds()

		return append(imageIds, snapshotIds...), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "The Image Alias to set instead of Image Id. A password or SSH Key need to be set")
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	cmd.AddStringFlag(cloudapiv6.ArgUserData, "", "", "The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	cmd.AddStringFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, "", "Absolute paths of the SSH Keys for the Volume")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume creation [seconds]")

	return cmd
}
