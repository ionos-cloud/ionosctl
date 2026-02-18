package image

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultImageCols = []string{"ImageId", "Name", "ImageAliases", "Location", "LicenceType", "ImageType", "CloudInit", "CreatedDate"}
	allImageCols     = []string{"ImageId", "Name", "ImageAliases", "Location", "Size", "LicenceType", "ImageType", "Description", "Public", "CloudInit", "CreatedDate", "CreatedBy", "CreatedByUserId", "ExposeSerial", "RequireLegacyBios", "ApplicationType"}
)

func ImageCmd() *core.Command {
	imageCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "image",
			Aliases:          []string{"img"},
			Short:            "Image Operations",
			Long:             "The sub-commands of `ionosctl image` allow you to see information about the Images available.",
			TraverseChildren: true,
		},
	}
	globalFlags := imageCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultImageCols, tabheaders.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetFlagName(imageCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = imageCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allImageCols, cobra.ShellCompDirectiveNoFileComp
	})

	imageCmd.AddCommand(ImageListCmd())
	imageCmd.AddCommand(ImageGetCmd())
	imageCmd.AddCommand(ImageUpdateCmd())
	imageCmd.AddCommand(ImageDeleteCmd())
	imageCmd.AddCommand(Upload())

	return core.WithConfigOverride(imageCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
