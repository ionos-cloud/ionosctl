package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func SnapshotUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Snapshot's metadata and inherited capabilities",
		LongDesc: `Use this command to update the metadata of an existing Snapshot. This edits properties recorded ON the Snapshot - it does not re-capture or change the stored data image.

Two kinds of properties can be changed:
  1. Descriptive: --name, --description, --licence-type, --sec-auth-protection.
  2. Hot-plug capability hints (cpu/ram/nic/disc SCSI/disc VirtIO, plug and unplug). These flags advertise whether a component can be added/removed at runtime without a reboot. A Volume created or restored from this Snapshot inherits these values, so setting them here fixes up the defaults future Volumes will get. Only pass the flags you want to change; unspecified capabilities are left untouched. VirtIO disk hot-unplug is unsupported on Windows, and SCSI disk hot-unplug is limited to non-Windows guests.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Snapshot Id`,
		Example: `# Rename a snapshot and fix its OS licence
ionosctl compute snapshot update --snapshot-id SNAPSHOT_ID --name "prod-db golden v2" --licence-type LINUX

# Advanced: mark the image as CPU/RAM/NIC hot-plug capable so volumes restored from it inherit those capabilities
ionosctl compute snapshot update --snapshot-id SNAPSHOT_ID --cpu-hot-plug=true --ram-hot-plug=true --nic-hot-plug=true --disc-virtio-plug=true --wait`,
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A human-friendly label for the Snapshot; shown in listings. Does not have to be unique")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Free-form notes about the Snapshot, e.g. why or when it was taken")
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "", constants.EnumLicenceType, "The operating-system licence recorded on the Snapshot. Inherited by Volumes created/restored from it and affects OS licensing (WINDOWS variants are billed)")
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "Advertise that Volumes from this Snapshot support adding vCPUs at runtime without a reboot. E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "Advertise that Volumes from this Snapshot support removing vCPUs at runtime without a reboot. E.g.: --cpu-hot-unplug=true, --cpu-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "Advertise that Volumes from this Snapshot support adding memory at runtime without a reboot. E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "Advertise that Volumes from this Snapshot support removing memory at runtime without a reboot. E.g.: --ram-hot-unplug=true, --ram-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "Advertise that Volumes from this Snapshot support attaching a NIC at runtime without a reboot. E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "Advertise that Volumes from this Snapshot support detaching a NIC at runtime without a reboot. E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "Advertise that Volumes from this Snapshot support attaching a VirtIO disk at runtime without a reboot. E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "Advertise that Volumes from this Snapshot support detaching a VirtIO disk at runtime without a reboot. Not supported on Windows guests. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", false, "Advertise that Volumes from this Snapshot support attaching a SCSI disk at runtime without a reboot. E.g.: --disc-scsi-plug=true, --disc-scsi-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "Advertise that Volumes from this Snapshot support detaching a SCSI disk at runtime without a reboot. Limited to non-Windows guests. E.g.: --disc-scsi-unplug=true, --disc-scsi-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgSecAuthProtection, "", false, "Protect the Snapshot with secure authentication: when true, deleting or restoring it requires the Contract Owner or a re-authenticated user. E.g.: --sec-auth-protection=true, --sec-auth-protection=false")

	return cmd
}
