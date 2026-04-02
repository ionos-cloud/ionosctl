package volume

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId",
		"DeviceNumber", "UserData", "BootServerId", "DatacenterId"}
)

func ServerVolumeCmd() *core.Command {
	serverVolumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Server Volume Operations",
			Long:             "The sub-commands of `ionosctl compute server volume` allow you to attach, get, list, detach Volumes from Servers.",
			TraverseChildren: true,
		},
	}

	serverVolumeCmd.AddCommand(ServerVolumeAttachCmd())
	serverVolumeCmd.AddCommand(ServerVolumeListCmd())
	serverVolumeCmd.AddCommand(ServerVolumeGetCmd())
	serverVolumeCmd.AddCommand(ServerVolumeDetachCmd())

	return core.WithConfigOverride(serverVolumeCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
