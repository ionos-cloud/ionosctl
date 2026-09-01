package image

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
)

func ImageUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "image",
		Resource:  "image",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a private image's properties and hardware capabilities",
		LongDesc: `Update the metadata and advertised hardware capabilities of one of your PRIVATE images (PUBLIC IONOS images cannot be modified). Typical uses are correcting the licence-type after an FTP upload, toggling cloud-init support, or declaring which resources can be hot-plugged.

The hot-plug / hot-unplug flags describe what a server booted from this image will support at runtime, i.e. whether the guest OS can have that resource added (hot-plug) or removed (hot-unplug) while the VM is running, without a reboot. Set them to match what your OS and drivers actually support; enabling a capability the guest cannot handle leads to failed or ignored operations.

Constraints:
  * You can only enable hot-UNPLUG for a resource whose hot-PLUG you also enabled.
  * Disk hot-unplug (--disc-virtio-hot-unplug, --disc-scsi-hot-unplug) is not supported for Windows guests.
  * --application-type only accepts a value other than UNKNOWN on PUBLIC images, so it is effectively a no-op on your private images.

Note: all boolean capability flags carry a default, so every one is sent on update; pass the flags explicitly to set the values you want.

Required values to run command:

* Image Id`,
		Example: `# Correct the licence type and name of an uploaded image
ionosctl compute image update --image-id IMAGE_ID --name "ubuntu-24.04-custom" --licence-type LINUX

# Declare a Linux image that supports CPU/RAM/NIC hot-plug but no disk hot-unplug, with cloud-init
ionosctl compute image update --image-id IMAGE_ID --licence-type LINUX --cloud-init V1 \
  --cpu-hot-plug=true --ram-hot-plug=true --nic-hot-plug=true \
  --disc-virtio-hot-unplug=false --disc-scsi-hot-unplug=false`,
		PreCmdRun:  PreRunImageId,
		CmdRun:     RunImageUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			return request.Filter("public", "false")
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end

	// Properties flags
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Human-friendly display name of the image (does not have to be unique)")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Free-text description of the image")
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", constants.EnumLicenceType, "OS licence type. Determines how IONOS bills and configures the guest (e.g. Windows editions are licensed). Use LINUX/RHEL for Linux, WINDOWS2016/2019/2022/2025 for the matching Windows Server, OTHER/UNKNOWN otherwise")
	cmd.AddSetFlag(constants.FlagCloudInit, "", "V1", []string{"V1", "NONE"}, "Whether servers built from this image accept cloud-init user-data for first-boot provisioning. V1 enables the cloud-init datasource; NONE disables it")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", true, "Guest supports adding CPU cores while running (no reboot)")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", true, "Guest supports adding RAM while running (no reboot)")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", true, "Guest supports attaching a NIC while running (no reboot)")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", true, "Guest supports attaching a Virt-IO disk while running (no reboot)")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", true, "Guest supports attaching a SCSI disk while running (no reboot)")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "Guest supports removing CPU cores while running. Only valid if CPU hot-plug is also enabled")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "Guest supports removing RAM while running. Only valid if RAM hot-plug is also enabled")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "Guest supports detaching a NIC while running. Only valid if NIC hot-plug is also enabled")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "Guest supports detaching a Virt-IO disk while running. Not supported on Windows guests")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "Guest supports detaching a SCSI disk while running. Not supported on Windows guests")
	cmd.AddBoolFlag(cloudapiv6.ArgExposeSerial, "", false, "Expose the attached disk's serial id to the guest. Some OSes/software need it; note it can influence licensed-software (e.g. Windows) behavior")
	cmd.AddBoolFlag(cloudapiv6.ArgRequireLegacyBios, "", true, "Boot the image in legacy BIOS mode instead of UEFI. Set false for images that require/expect UEFI")
	cmd.AddSetFlag(cloudapiv6.ArgApplicationType, "", "UNKNOWN", constants.EnumApplicationType, "Application pre-installed on the image (e.g. an MSSQL edition). Only PUBLIC images may set a value other than UNKNOWN, so this is a no-op on your private images")

	return cmd
}
