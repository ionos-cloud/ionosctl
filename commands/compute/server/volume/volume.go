package volume

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allVolumeCols = []table.Column{
	{Name: "VolumeId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "Type", JSONPath: "properties.type", Default: true},
	{Name: "LicenceType", JSONPath: "properties.licenceType", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Image", JSONPath: "properties.image", Default: true},
	{Name: "Bus", JSONPath: "properties.bus"},
	{Name: "AvailabilityZone", JSONPath: "properties.availabilityZone"},
	{Name: "BackupunitId", JSONPath: "properties.backupunitId"},
	{Name: "DeviceNumber", JSONPath: "properties.deviceNumber"},
	{Name: "UserData", JSONPath: "properties.userData"},
	{Name: "BootServerId", JSONPath: "properties.bootServer"},
	{Name: "DatacenterId", JSONPath: "href"},
}

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
