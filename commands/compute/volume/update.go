package volume

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VolumeUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Volume",
		LongDesc: `Update the mutable properties of an existing Volume.

Resizing: --size may only GROW the Volume; the IONOS CLOUD API cannot shrink a Volume once provisioned. If the attached Server (and the guest OS) supports disk hot-plug, the new capacity appears live without a reboot. The extra space is raw - it is NOT added to any partition or filesystem automatically, so you must extend the partition/filesystem from inside the operating system afterwards.

Immutable properties: the storage tier (--type), availability zone and the bootable image/credentials are fixed at creation and cannot be changed here. --name and --bus can be adjusted; the hot-plug capability flags advertise what the disk supports to the guest.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the Volume to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example: `# Grow a volume to 20 GB (extend the filesystem inside the OS afterwards)
ionosctl compute volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --size 20GB

# Rename a volume and enable RAM hot-plug advertisement
ionosctl compute volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name prod-data --ram-hot-plug=true`,
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunVolumeUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A new human-friendly label for the Volume")
	cmd.AddStringFlag(cloudapiv6.ArgSize, "", "", "The new capacity of the Volume. Can only be increased, never decreased. Accepts a plain number (GB) or a unit suffix, e.g. --size 20 or --size 20GB. Upper bound is 4 TB (larger on request) and your contract limit. Remember to extend the partition/filesystem inside the guest OS afterwards")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "The virtual bus the disk is exposed on. VIRTIO is the high-performance default; IDE is a legacy fallback")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "Advertise to the guest OS that CPUs can be added without a reboot. E.g.: --cpu-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "Advertise to the guest OS that memory can be added without a reboot. E.g.: --ram-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "Advertise to the guest OS that a NIC can be added without a reboot. E.g.: --nic-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "Advertise to the guest OS that a NIC can be removed without a reboot. E.g.: --nic-hot-unplug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "Advertise to the guest OS that a VirtIO storage volume can be attached without a reboot. E.g.: --disc-virtio-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "Advertise to the guest OS that a VirtIO storage volume can be detached without a reboot. Not supported by Windows guests. E.g.: --disc-virtio-hot-unplug=true")

	return cmd
}
