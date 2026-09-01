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

func ImageDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "image",
		Resource:  "image",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a private image",
		LongDesc: `Delete one of your PRIVATE images, or all of them with --all. PUBLIC (IONOS-provided) images cannot be deleted and are always skipped.

IMPORTANT: for an image you uploaded via FTP, this API call only sets the image size to 0B — it does NOT remove the underlying file from the FTP server. To fully remove an FTP-uploaded image you must contact IONOS support.

Required values to run command:

* Image Id`,
		Example:    "ionosctl compute image delete --image-id IMAGE_ID",
		PreCmdRun:  PreRunImageDelete,
		CmdRun:     RunImageDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			return request.Filter("public", "false")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete every PRIVATE (non-public) image on the contract. Public images are skipped because the API forbids deleting them")

	return cmd
}
