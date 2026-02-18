package volume

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId",
		"DeviceNumber", "UserData", "BootServerId", "DatacenterId"}
)

func VolumeCmd() *core.Command {
	volumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Volume Operations",
			Long:             "The sub-commands of `ionosctl volume` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes. To attach a Volume to a Server, use the Server command `ionosctl server volume attach`.",
			TraverseChildren: true,
		},
	}
	globalFlags := volumeCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = viper.BindPFlag(core.GetFlagName(volumeCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})

	volumeCmd.AddCommand(VolumeListCmd())
	volumeCmd.AddCommand(VolumeGetCmd())
	volumeCmd.AddCommand(VolumeCreateCmd())
	volumeCmd.AddCommand(VolumeUpdateCmd())
	volumeCmd.AddCommand(VolumeDeleteCmd())

	return core.WithConfigOverride(volumeCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
