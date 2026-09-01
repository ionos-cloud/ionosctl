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

func VolumeCmd() *core.Command {
	volumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "volume",
			Aliases: []string{"v", "vol"},
			Short:   "Volume Operations",
			Long: `Manage block storage volumes inside a Virtual Data Center (VDC).

A Volume is a virtual block device that lives in one datacenter (--datacenter-id) and, once created, can be attached to a Server to act as its disk. Volumes are provisioned and billed independently of servers: creating a Volume does NOT attach it. Use ` + "`" + `ionosctl compute server volume attach` + "`" + ` to connect a Volume to a Server, at which point the guest OS sees it as a disk on the chosen bus (VIRTIO or IDE).

Volume model:
  - Storage tier (--type): HDD, SSD Standard, SSD Premium, or DAS (Direct-Attached Storage on Cube instances); the performance classes ESSENTIAL, BALANCED and PERFORMANCE (performance-tiered volumes) are also accepted. The tier fixes the price/performance profile and cannot be changed after provisioning.
  - Capacity (--size): can be grown later, never shrunk. A boot volume additionally carries an OS image plus its initial credentials.
  - Bootability: a Volume becomes bootable by being created from an image or snapshot (--image / --image-alias). A blank Volume with only a --licence-type is a raw data disk.
  - Location: the storage availability zone (--availability-zone) is chosen at creation and is immutable afterwards.`,
			TraverseChildren: true,
		},
	}
	volumeCmd.AddColsFlag(allVolumeCols)

	volumeCmd.AddCommand(VolumeListCmd())
	volumeCmd.AddCommand(VolumeGetCmd())
	volumeCmd.AddCommand(VolumeCreateCmd())
	volumeCmd.AddCommand(VolumeUpdateCmd())
	volumeCmd.AddCommand(VolumeDeleteCmd())

	return core.WithConfigOverride(volumeCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
