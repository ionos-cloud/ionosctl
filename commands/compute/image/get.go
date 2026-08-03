package image

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ImageGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "image",
		Resource:  "image",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a specified Image",
		LongDesc: `Retrieve the full properties of a single Image by its UUID: its licence-type, location, size, image-aliases, cloud-init support, hot-plug capabilities and (for Confidential Computing images) the SEV-SNP entry in RequiredFeatures. Works for both PUBLIC and PRIVATE images.

Required values to run command:

* Image Id`,
		Example:    getImageExample,
		PreCmdRun:  PreRunImageId,
		CmdRun:     RunImageGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
