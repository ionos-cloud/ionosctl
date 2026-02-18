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

func ServerCreateCmd() *core.Command {
	create := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Server",
		LongDesc: `Use this command to create an ENTERPRISE, CUBE, VCPU or GPU Server in a specified Virtual Data Center.

1. For ENTERPRISE Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

To see which CPU Family are available in which location, use ` + "`" + `ionosctl location` + "`" + ` commands.

Required values to create a Server of type ENTERPRISE:

* Data Center Id
* Cores
* RAM

2. For CUBE Servers:

Servers of type CUBE will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use ` + "`" + `ionosctl template` + "`" + ` commands.

Required values to create a Server of type CUBE:

* Data Center Id
* Type
* Template Id

3. For VCPU Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

You cannot set the CPU Family for VCPU Servers.

Required values to create a Server of type VCPU:

* Data Center Id
* Type
* Cores
* RAM

4. For GPU Servers:

Servers of type GPU will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use ` + "`" + `ionosctl template` + "`" + ` commands.

GPU servers do not support the --cpu-family flag and are automatically assigned the AMD_TURIN CPU family.

Required values to create a Server of type GPU:

* Data Center Id
* Type
* Template Id

By default, Licence Type for Direct Attached Storage is set to LINUX. You can set it using the ` + "`" + `--licence-type` + "`" + ` option or set an Image Id. For Image Id, it is needed to set a password or SSH keys.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.`,
		Example:    "ionosctl server create --datacenter-id DATACENTER_ID --cores 2 --ram 256MB\nionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID\nionosctl server create --datacenter-id DATACENTER_ID --type VCPU --cores 2 --ram 256MB\nionosctl server create --datacenter-id DATACENTER_ID --type GPU --template-id TEMPLATE_ID",
		PreCmdRun:  PreRunServerCreate,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Server", "Name of the Server")
	create.AddIntFlag(constants.FlagCores, "", cloudapiv6.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCpuFamily, "", cloudapiv6.DefaultServerCPUFamily,
		"CPU Family for the Server. For CUBE Servers, the CPU Family is INTEL_SKYLAKE. "+
			"If the flag is not set, the CPU Family will be chosen based on the location of the Datacenter. "+
			"It will always be the first CPU Family available, as returned by the API")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(create.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagNICMultiQueue, "", false, constants.FlagNICMultiQueueDescription)
	create.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgTemplateId, "", "", "[CUBE Server] The unique Template Id", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagType, "", "ENTERPRISE", []string{"ENTERPRISE", "CUBE", "VCPU", "GPU"}, "Type usages for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ENTERPRISE", "CUBE", "VCPU", "GPU"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagPromoteVolume, "", false, "For CUBE and GPU servers, promotes the attached volume to be the Boot Volume. Requires --wait-for-state")

	// Volume Properties - for DAS Volume associated with Cube Server
	create.AddStringFlag(cloudapiv6.ArgVolumeName, "N", "Unnamed Direct Attached Storage", "[CUBE Server] Name of the Direct Attached Storage")
	create.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "[CUBE Server] The bus type of the Direct Attached Storage")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(cloudapiv6.ArgLicenceType, "l", "LINUX", constants.EnumLicenceType, "[CUBE Server] Licence Type of the Direct Attached Storage")
	create.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "[CUBE Server] The Image Alias to use instead of Image Id for the Direct Attached Storage")
	create.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "[CUBE Server] The Image Id or snapshot Id to be used as for the Direct Attached Storage")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		imageIds := completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "HDD")
		})

		snapshotIds := completer.SnapshotIds()

		return append(imageIds, snapshotIds...), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "[CUBE Server] Initial image password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	create.AddStringSliceFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, []string{""}, "[CUBE Server] Absolute paths for the SSH Keys of the Direct Attached Storage")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server creation to be executed")
	create.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for new Server to be in AVAILABLE state")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds]")

	return create
}
