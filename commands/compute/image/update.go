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
		Namespace:  "image",
		Resource:   "image",
		Verb:       "update",
		Aliases:    []string{"u", "up"},
		ShortDesc:  "Update a specified Image",
		LongDesc:   "Use this command to update information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    "ionosctl compute image update --image-id IMAGE_ID",
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

	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")

	// Properties flags
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Image")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Image")
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", constants.EnumLicenceType, "The OS type of this image")
	cmd.AddSetFlag(constants.FlagCloudInit, "", "V1", []string{"V1", "NONE"}, "Cloud init compatibility")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", true, "'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", true, "'Hot-Plug' RAM")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", true, "'Hot-Plug' NIC")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", true, "'Hot-Plug' Virt-IO drive")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", true, "'Hot-Plug' SCSI drive")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "'Hot-Unplug' RAM")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "'Hot-Unplug' NIC")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "'Hot-Unplug' Virt-IO drive")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "'Hot-Unplug' SCSI drive")
	cmd.AddBoolFlag(cloudapiv6.ArgExposeSerial, "", false, "If set to `true` will expose the serial id of the disk attached to the server")
	cmd.AddBoolFlag(cloudapiv6.ArgRequireLegacyBios, "", true, "Indicates if the image requires the legacy BIOS for compatibility or specific needs.")
	cmd.AddSetFlag(cloudapiv6.ArgApplicationType, "", "UNKNOWN", constants.EnumApplicationType, "The type of application that is hosted on this resource")

	return cmd
}
