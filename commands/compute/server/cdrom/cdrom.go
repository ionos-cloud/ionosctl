package cdrom

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

func ServerCdromCmd() *core.Command {
	serverCdromCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cdrom",
			Aliases:          []string{"cd"},
			Short:            "Server CD-ROM Operations",
			Long:             "The sub-commands of `ionosctl server cdrom` allow you to attach, get, list, detach CD-ROMs from Servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := serverCdromCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultImageCols, tabheaders.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetFlagName(serverCdromCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = serverCdromCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allImageCols, cobra.ShellCompDirectiveNoFileComp
	})

	serverCdromCmd.AddCommand(ServerCdromAttachCmd())
	serverCdromCmd.AddCommand(ServerCdromListCmd())
	serverCdromCmd.AddCommand(ServerCdromGetCmd())
	serverCdromCmd.AddCommand(ServerCdromDetachCmd())

	return core.WithConfigOverride(serverCdromCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
