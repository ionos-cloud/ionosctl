package volume

import (
	"context"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
)

func VolumeCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Volume",
		LongDesc: `Create a block storage Volume inside a Virtual Data Center. This does NOT attach the Volume to a Server; attach it afterwards with ` + "`" + `ionosctl compute server volume attach` + "`" + `.

Storage tier (--type) determines the price/performance profile and is fixed for the life of the Volume:
  * HDD          - spinning disks, lowest cost. Best for backups, archives and cold storage. (~1,100 IOPS, up to 2,500 burst.)
  * SSD Standard - general-purpose flash. ~40 read / 30 write IOPS per GB, up to 24,000/18,000 IOPS per volume.
  * SSD Premium  - high-performance flash for databases and latency-sensitive workloads. ~75 read / 50 write IOPS per GB, up to 45,000/30,000 IOPS per volume.
  * DAS          - Direct-Attached NVMe storage that ships with Cube instances. It is bound to the Cube, its size is fixed by the Cube template, and it cannot be resized, detached or deleted independently.
--type also accepts the performance classes ESSENTIAL, BALANCED and PERFORMANCE (performance-tiered volumes) as valid values.
Performance scales with volume size (IOPS-per-GB), so IONOS recommends booking SSD volumes of at least 100 GB. Volume sizes range from 1 GB up to 4 TB (larger on request).

Blank vs. bootable Volume:
  * Blank data disk: omit --image/--image-alias and set --licence-type so the platform knows how to bill/handle the disk. The disk is unformatted; partition and format it from the OS after attaching.
  * Bootable OS disk: pass --image (Image or Snapshot Id) OR --image-alias. When an image is set, --licence-type is derived automatically from the image and should not be overridden. For IONOS public images you must seed initial credentials with --password and/or --ssh-key-paths, otherwise you will not be able to log in. Setting --password even alongside SSH keys is recommended so the DCD remote console can authenticate with a password.

cloud-init: --user-data injects a base64-encoded cloud-init configuration on first boot (users, packages, scripts). It requires a cloud-init capable public image or image-alias.

Immutability: --type and --availability-zone are chosen at creation and cannot be changed afterwards. --size can be increased later but never decreased.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the Volume to reach AVAILABLE state before the command returns.

Required values to run command:

* Data Center Id`,
		Example: `# Blank 50GB SSD Premium data disk (attach to a server separately)
ionosctl compute volume create --datacenter-id DATACENTER_ID --name data-disk --size 50GB --type "SSD Premium" --licence-type LINUX

# Bootable Linux volume from a public image alias, seeded with SSH keys and a console password
ionosctl compute volume create --datacenter-id DATACENTER_ID --name boot-disk --size 20GB --type "SSD Standard" --image-alias ubuntu:latest --ssh-key-paths "$HOME/.ssh/id_rsa.pub" --password 'S3curePassw0rd!' --user-data "$(base64 -w0 cloud-init.yaml)"`,
		PreCmdRun:  PreRunVolumeCreate,
		CmdRun:     RunVolumeCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Volume", "A human-friendly label for the Volume. Not required to be unique")
	cmd.AddStringFlag(cloudapiv6.ArgSize, cloudapiv6.ArgSizeShort, strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The capacity of the Volume. Accepts a plain number (interpreted as GB) or a unit suffix, e.g. --size 10 or --size 10GB or --size 1TB. Range is 1 GB to 4 TB (larger on request); the upper bound is also capped by your contract limit. Can be increased later but never decreased")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "The virtual bus the disk is exposed on once attached to a Server. VIRTIO is the paravirtualized, high-performance default and is recommended for all modern OSes. IDE is a legacy, lower-performance bus needed only in special cases (e.g. temporarily during a Windows install before VirtIO drivers are available)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "LINUX", constants.EnumLicenceType, "The OS licence the Volume is billed and configured for. Use for blank data disks so the platform knows how to handle them. When --image or --image-alias is set, the licence type is derived automatically from the image and this flag should not be set")
	cmd.AddStringFlag(constants.FlagType, "", "HDD", "The storage tier (fixed for the life of the Volume). HDD is cheapest (backups/cold storage); 'SSD Standard' is general-purpose flash; 'SSD Premium' is high-IOPS flash for databases; DAS is the fixed NVMe disk bundled with Cube instances. The performance classes ESSENTIAL, BALANCED and PERFORMANCE (performance-tiered volumes) are also accepted")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ESSENTIAL", "BALANCED", "PERFORMANCE", "HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "The storage availability zone the Volume is physically placed in. AUTO lets the platform pick; ZONE_1/ZONE_2/ZONE_3 pin it to a specific zone (useful to spread replicas across failure domains). Immutable after provisioning - to move a Volume to another zone you must snapshot it and re-create from the snapshot")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, "", "", "The Id of a Backup Unit you own, used to schedule automatic backups of this Volume. Only valid on a bootable Volume, so it must be combined with a public --image or --image-alias")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "Id of an Image or Snapshot to clone onto the Volume, making it bootable. Mutually exclusive with --image-alias. For an IONOS public image you must also seed credentials with --password and/or --ssh-key-paths. The image must match the Volume's location and disk type")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		imageIds := completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			dcId, _ := c.Flags().GetString(cloudapiv6.ArgDataCenterId)
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				dcId).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return r
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "HDD")
		})

		snapshotIds := completer.SnapshotIds()

		return append(imageIds, snapshotIds...), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "A stable, human-readable alias for an IONOS public image (e.g. ubuntu:latest) to use instead of a raw --image Id. Same credential rules apply: seed --password and/or --ssh-key-paths")
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Initial root/Administrator password baked into the OS on first boot. Public images only, and immutable afterwards (rotate it from inside the guest). Allowed characters: a-z, A-Z, 0-9, 8-50 chars. Recommended even when using SSH keys so the DCD remote console can log in")
	cmd.AddStringFlag(cloudapiv6.ArgUserData, "", "", "A base64-encoded cloud-init configuration applied on first boot (create users, install packages, run scripts). Requires a cloud-init capable public --image or --image-alias. Encode a file with e.g. base64 -w0 cloud-init.yaml")
	cmd.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "Advertise to the guest OS that CPUs can be added to the attached Server without a reboot. E.g.: --cpu-hot-plug=true. Usually inherited from the image and rarely set by hand")
	cmd.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "Advertise to the guest OS that memory can be added to the attached Server without a reboot. E.g.: --ram-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "Advertise to the guest OS that a NIC can be added without a reboot. E.g.: --nic-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "Advertise to the guest OS that a NIC can be removed without a reboot. E.g.: --nic-hot-unplug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "Advertise to the guest OS that a VirtIO storage volume can be attached without a reboot. E.g.: --disc-virtio-hot-plug=true")
	cmd.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "Advertise to the guest OS that a VirtIO storage volume can be detached without a reboot. Not supported by Windows guests. E.g.: --disc-virtio-hot-unplug=true")
	cmd.AddStringFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, "", "Comma-separated absolute paths to public SSH key files to authorize for the image's default user on first boot. Public images only. e.g. --ssh-key-paths \"$HOME/.ssh/id_rsa.pub,/keys/ops.pub\"")

	return cmd
}
