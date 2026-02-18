package snapshot

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
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
		ShortDesc: "Update a Snapshot",
		LongDesc: `Use this command to update a specified Snapshot.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Snapshot Id`,
		Example:    "ionosctl snapshot update --snapshot-id SNAPSHOT_ID --name NAME",
		PreCmdRun:  cloudapiv6cmds.PreRunSnapshotId,
		CmdRun:     cloudapiv6cmds.RunSnapshotUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Snapshot")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Snapshot")
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "", constants.EnumLicenceType, "Licence Type of the Snapshot")
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "This volume is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "This volume is capable of CPU hot unplug (no reboot required). E.g.: --cpu-hot-unplug=true, --cpu-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "This volume is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "This volume is capable of memory hot unplug (no reboot required). E.g.: --ram-hot-unplug=true, --ram-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "This volume is capable of NIC hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "This volume is capable of NIC hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "This volume is capable of VirtIO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "This volume is capable of VirtIO drive hot unplug (no reboot required). E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", false, "This volume is capable of SCSI drive hot plug (no reboot required). E.g.: --disc-scsi-plug=true, --disc-scsi-plug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "This volume is capable of SCSI drive hot unplug (no reboot required). E.g.: --disc-scsi-unplug=true, --disc-scsi-unplug=false")
	cmd.AddBoolFlag(cloudapiv6.ArgSecAuthProtection, "", false, "Enable secure authentication protection. E.g.: --sec-auth-protection=true, --sec-auth-protection=false")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Snapshot creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Snapshot creation [seconds]")

	return cmd
}
