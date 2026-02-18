package image

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	cloudapiv6image "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/image"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
)

func ImageDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete an image",
		LongDesc:   "Use this command to delete a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    "ionosctl image delete --image-id IMAGE_ID",
		PreCmdRun:  cloudapiv6image.PreRunImageDelete,
		CmdRun:     cloudapiv6image.RunImageDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			return request.Filter("public", "false")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all non-public images")

	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")

	return cmd
}
