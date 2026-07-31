package server

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerUpdateCmd() *core.Command {
	update := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update the compute shape or boot configuration of a Server",
		LongDesc: `Use this command to update a Server in a Virtual Data Center. You can rescale an ENTERPRISE/VCPU server (--cores, --ram), change its --cpu-family or --availability-zone, and repoint what it boots from (--volume-id for a boot Volume, --cdrom-id for a boot CD-ROM). Both boot targets must already be attached to the Server.

RAM (--ram) must be a multiple of 256. The default unit is MB, so --ram 256 = 256MB; a unit may be given, e.g. --ram 1GB. Minimum 256MB; the maximum depends on your contract limit.

CUBE constraint: for CUBE Servers only the Name can be updated — their cores, RAM and CPU family are fixed by the instance-size template and cannot be changed here. Some changes (e.g. --nic-multi-queue, CPU family) require a server restart to take effect.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id`,
		Example: `# BASIC: rescale an ENTERPRISE server to 4 cores
ionosctl compute server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 4

# ADVANCED: rescale cores + RAM and set the boot volume to an already-attached Volume, waiting until AVAILABLE
ionosctl compute server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 8 --ram 16GB --volume-id VOLUME_ID --wait`,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", "Id of the Volume to set as the Server's boot volume. The Volume must already be attached to this Server. Mutually exclusive with booting from a CD-ROM (--cdrom-id)")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgCdromId, "", "", "Id of the CD-ROM to set as the Server's boot device. The CD-ROM must already be attached to this Server. Use this to boot from an installer ISO instead of a Volume (--volume-id)")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for CDROM images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "CDROM")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New display name for the Server. This is the only property updatable on CUBE Servers")
	update.AddStringFlag(constants.FlagCpuFamily, "", "", "New CPU family for the Server, e.g. INTEL_SKYLAKE, INTEL_XEON, AMD_OPTERON. Availability depends on the datacenter location; changing it requires a server restart. Not applicable to CUBE/VCPU/GPU")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(update.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.FlagNICMultiQueue, "", false, constants.FlagNICMultiQueueDescription)
	update.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "", "Physical availability zone of the Server: AUTO (platform-chosen), ZONE_1 or ZONE_2")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(constants.FlagCores, "", cloudapiv6.DefaultServerCores, "New number of CPU cores (ENTERPRISE/VCPU only; fixed by the template for CUBE/GPU). Maximum depends on your contract resource limits")
	update.AddStringFlag(constants.FlagRam, "", "", "New memory size, in multiples of 256. Default unit is MB (e.g. --ram 256 = 256MB); a unit may be given (e.g. --ram 1GB). Minimum 256MB, maximum per contract limit. ENTERPRISE/VCPU only")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})

	return update
}
