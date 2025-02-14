package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	serverCubeType       = "CUBE"
	serverEnterpriseType = "ENTERPRISE"
	serverVCPUType       = "VCPU"
)

var (
	DefaultServerCols = []string{"ServerId", "Name", "Type", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State"}
	AllServerCols     = []string{"ServerId", "DatacenterId", "Name", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State", "TemplateId", "Type", "BootCdromId", "BootVolumeId"}
)

func ServerCmd() *core.Command {
	ctx := context.TODO()
	serverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "svr"},
			Short:            "Server Operations",
			Long:             "The sub-commands of `ionosctl server` allow you to manage Servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", DefaultServerCols, tabheaders.ColsMessage(AllServerCols))
	_ = viper.BindPFlag(core.GetFlagName(serverCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = serverCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return AllServerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "server",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Servers",
		LongDesc:   "Use this command to list Servers from a specified Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ServersFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listServerExample,
		PreCmdRun:  PreRunServerList,
		CmdRun:     RunServerList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "server",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Server",
		LongDesc:   "Use this command to get information about a specified Server from a Virtual Data Center. You can also wait for Server to get in AVAILABLE state using `--wait-for-state` option.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    getServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for specified Server to be in AVAILABLE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for waiting for Server to be in AVAILABLE state [seconds]")
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Server",
		LongDesc: `Use this command to create an ENTERPRISE, CUBE or VCPU Server in a specified Virtual Data Center.

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
* Cores
* RAM

By default, Licence Type for Direct Attached Storage is set to LINUX. You can set it using the ` + "`" + `--licence-type` + "`" + ` option or set an Image Id. For Image Id, it is needed to set a password or SSH keys.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.`,
		Example:    createServerExample,
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
	create.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgTemplateId, "", "", "[CUBE Server] The unique Template Id", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagType, "", serverEnterpriseType, []string{serverEnterpriseType, serverCubeType, serverVCPUType}, "Type usages for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{serverEnterpriseType, serverCubeType}, cobra.ShellCompDirectiveNoFileComp
	})

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
		imageIds := completer.ImageIds(func(r compute.ApiImagesGetRequest) compute.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return compute.ApiImagesGetRequest{}
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
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Server",
		LongDesc: `Use this command to update a specified Server from a Virtual Data Center.

You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit.

Note: For CUBE Servers, only Name attribute can be updated.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    updateServerExample,
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
	update.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", "The unique Volume Id for the BootVolume. The Volume needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgCdromId, "", "", "The unique Cdrom Id for the BootCdrom. The Cdrom needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(r compute.ApiImagesGetRequest) compute.ApiImagesGetRequest {
			// Completer for CDROM images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return compute.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "CDROM")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Server")
	update.AddStringFlag(constants.FlagCpuFamily, "", "", "CPU Family of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagCpuFamily, func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		datacenterId := viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))
		return completer.DatacenterCPUFamilies(update.Command.Context(), datacenterId), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "", "Availability zone of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(constants.FlagCores, "", cloudapiv6.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server update to be executed")
	update.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for the updated Server to be in AVAILABLE state")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server update/for Server to be in AVAILABLE state [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Server",
		LongDesc: `Use this command to delete a specified Server from a Virtual Data Center.

NOTE: This will not automatically remove the storage Volumes attached to a Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    deleteServerExample,
		PreCmdRun:  PreRunDcServerDelete,
		CmdRun:     RunServerDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Servers form a virtual Datacenter.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	/*
		Suspend Command
	*/
	suspend := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "suspend",
		ShortDesc: "Suspend a Cube Server",
		LongDesc: `Use this command to suspend a Cube Server. The operation can only be applied to Cube Servers. Note: The virtual machine will not be deleted.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    suspendServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerSuspend,
		InitClient: true,
	})
	suspend.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(viper.GetString(core.GetFlagName(suspend.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string][]string{
					"type": {serverCubeType},
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server suspend to be executed")
	suspend.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server suspend [seconds]")
	suspend.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Start Command
	*/
	start := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "start",
		Aliases:   []string{"on"},
		ShortDesc: "Start a Server",
		LongDesc: `Use this command to start a Server from a Virtual Data Center. If the Server's public IP was deallocated then a new IP will be assigned.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    startServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerStart,
		InitClient: true,
	})
	start.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(start.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server start to be executed")
	start.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server start [seconds]")
	start.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Stop Command
	*/
	stop := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "stop",
		Aliases:   []string{"off"},
		ShortDesc: "Stop a Server",
		LongDesc: `Use this command to stop a Server from a Virtual Data Center. The machine will be forcefully powered off, billing will cease, and the public IP, if one is allocated, will be deallocated.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    stopServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerStop,
		InitClient: true,
	})
	stop.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(stop.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server stop to be executed")
	stop.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server stop [seconds]")
	stop.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Reboot Command
	*/
	reboot := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "reboot",
		Aliases:   []string{"r"},
		ShortDesc: "Force a hard reboot of a Server",
		LongDesc: `Use this command to force a hard reboot of the Server. Do not use this method if you want to gracefully reboot the machine. This is the equivalent of powering off the machine and turning it back on.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    rebootServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerReboot,
		InitClient: true,
	})
	reboot.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(reboot.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server reboot to be executed")
	reboot.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")
	reboot.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Resume Command
	*/
	resume := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "resume",
		Aliases:   []string{"res"},
		ShortDesc: "Resume a Cube Server",
		LongDesc: `Use this command to resume a Cube Server. The operation can only be applied to suspended Cube Servers.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    resumeServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerResume,
		InitClient: true,
	})
	resume.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddUUIDFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(viper.GetString(core.GetFlagName(resume.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string][]string{
					"type": {serverCubeType},
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Server resume to be executed")
	resume.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Server resume [seconds]")
	resume.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	serverCmd.AddCommand(ServerTokenCmd())
	serverCmd.AddCommand(ServerConsoleCmd())
	serverCmd.AddCommand(ServerVolumeCmd())
	serverCmd.AddCommand(ServerCdromCmd())

	return serverCmd
}

func PreRunServerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ServersFilters(), completer.ServersFiltersUsage())
	}
	return nil
}

func PreRunServerCreate(c *core.PreCommandConfig) error {
	baseRequiredFlags := [][]string{
		{cloudapiv6.ArgDataCenterId, constants.FlagCores, constants.FlagRam},
		{cloudapiv6.ArgDataCenterId, constants.FlagType, cloudapiv6.ArgTemplateId},
	}

	// Start by checking the base required flags.
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, baseRequiredFlags...)
	if err != nil {
		return err
	}

	imageIdFlag := core.GetFlagName(c.NS, cloudapiv6.ArgImageId)
	imageAliasFlag := core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)

	// Check if image ID or alias is set
	if viper.IsSet(imageIdFlag) || viper.IsSet(imageAliasFlag) {
		imageRequiredFlags := make([][]string, 0)

		if viper.IsSet(imageAliasFlag) {
			// Handle public image alias
			for _, baseFlagSet := range baseRequiredFlags {
				imageRequiredFlags = append(imageRequiredFlags,
					append(baseFlagSet, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword),
					append(baseFlagSet, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths),
				)
			}
		} else if viper.IsSet(imageIdFlag) {
			// Check if the image ID corresponds to an image or snapshot
			img, _, imgErr := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
				viper.GetString(imageIdFlag)).Execute()
			if imgErr != nil {
				// Try to fetch as a snapshot if image fetch fails
				_, _, snapshotErr := client.Must().CloudClient.SnapshotsApi.SnapshotsFindById(context.Background(),
					viper.GetString(imageIdFlag)).Execute()
				if snapshotErr != nil {
					return fmt.Errorf("failed getting image or snapshot %s: %w", viper.GetString(imageIdFlag), imgErr)
				}

				// If it's a snapshot, no additional checks are required
				return nil
			}

			// If it's an image, determine if it is public or private
			for _, baseFlagSet := range baseRequiredFlags {
				if img.Properties != nil && img.Properties.Public != nil && *img.Properties.Public {
					// For public images, require password or SSH key
					imageRequiredFlags = append(imageRequiredFlags,
						append(baseFlagSet, cloudapiv6.ArgImageId, cloudapiv6.ArgPassword),
						append(baseFlagSet, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths),
					)
				} else {
					// For private images, only the image ID is required
					imageRequiredFlags = append(imageRequiredFlags,
						append(baseFlagSet, cloudapiv6.ArgImageId),
					)
				}
			}
		}

		err = core.CheckRequiredFlagsSets(c.Command, c.NS, imageRequiredFlags...)
		if err != nil {
			return err
		}
	}

	return nil
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)
}

func PreRunDcServerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunServerListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters

			if val, ok := filters["ram"]; ok {
				convertedSize, err := utils2.ConvertSize(val[0], utils2.MegaBytes)
				if err != nil {
					return err
				}

				filters["ram"] = []string{strconv.Itoa(convertedSize)}
				listQueryParams.Filters = &filters
			}
		}
	}

	// Don't apply listQueryParams to parent resource, as it would have unexpected side effects on the results
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	allDcs := getDataCenters(datacenters)
	var allServers []compute.Servers
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		servers, resp, err := c.CloudApiV6Services.Servers().List(*id, listQueryParams)
		if err != nil {
			return err
		}

		allServers = append(allServers, servers.Servers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.Server, allServers, tabheaders.GetHeaders(AllServerCols, DefaultServerCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunServerListAll(c)
	}
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters

			if val, ok := filters["ram"]; ok {
				convertedSize, err := utils2.ConvertSize(val[0], utils2.MegaBytes)
				if err != nil {
					return err
				}

				filters["ram"] = []string{strconv.Itoa(convertedSize)}
				listQueryParams.Filters = &filters
			}
		}
	}

	servers, resp, err := c.CloudApiV6Services.Servers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Server, servers.Servers,
		tabheaders.GetHeaders(AllServerCols, DefaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Server with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))))

	if err := waitfor.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
		return err
	}
	svr, resp, err := c.CloudApiV6Services.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(AllServerCols, DefaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	input, err := getNewServer(c)
	if err != nil {
		return err
	}

	// If Server is of type CUBE, it will create an attached Volume
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverCubeType {
		// Volume Properties
		volumeDAS, err := getNewDAS(c)
		if err != nil {
			return err
		}

		// Attach Storage
		input.SetEntities(compute.ServerEntities{
			Volumes: &compute.AttachedVolumes{
				Items: &[]compute.Volume{volumeDAS.Volume},
			},
		})
	}

	svr, resp, err := c.CloudApiV6Services.Servers().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := svr.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.ServerStateInterrogator, *id); err != nil {
				return err
			}

			if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
				*id, queryParams); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new server id")
		}
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(AllServerCols, DefaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	input, err := getUpdateServerInfo(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Updating Server with ID: %v in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))

	svr, resp, err := c.CloudApiV6Services.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		*input,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err = waitfor.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
			return err
		}

		if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)), queryParams); err != nil {
			return err
		}
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Server, svr.Server,
		tabheaders.GetHeaders(AllServerCols, DefaultServerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllServers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Server with id: %v from datacenter with id: %v... ", serverId, dcId))

	resp, err := c.CloudApiV6Services.Servers().Delete(dcId, serverId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully deleted"))
	return nil
}

func RunServerStart(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "start server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server is starting... "))

	resp, err := c.CloudApiV6Services.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully started"))
	return nil
}

func RunServerStop(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "stop server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server is stopping... "))

	resp, err := c.CloudApiV6Services.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully stopped"))
	return nil
}

func RunServerSuspend(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "suspend cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server is Suspending... "))

	resp, err := c.CloudApiV6Services.Servers().Suspend(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully suspended"))
	return nil
}

func RunServerReboot(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "reboot server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server is rebooting... "))

	resp, err := c.CloudApiV6Services.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully rebooted"))
	return nil
}

func RunServerResume(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "resume cube server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server is resuming... "))

	resp, err := c.CloudApiV6Services.Servers().Resume(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server successfully resumed"))
	return nil
}

func getUpdateServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := compute.ServerProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property name set: %v ", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCpuFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily))
		input.SetCpuFamily(cpuFamily)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CpuFamily set: %v ", cpuFamily))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
		input.SetAvailabilityZone(availabilityZone)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property AvailabilityZone set: %v ", availabilityZone))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		input.SetCores(cores)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Cores set: %v ", cores))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)) {
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
		input.SetBootVolume(compute.ResourceReference{
			Id: &volumeId,
		})

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property BootVolume set: %v ", volumeId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)) {
		cdromId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId))
		input.SetBootCdrom(compute.ResourceReference{
			Id: &cdromId,
		})

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property BootCdrom set: %v ", cdromId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		size, err := utils2.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
			utils2.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		input.SetRam(int32(size))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB ", int32(size)))
	}

	return &resources.ServerProperties{
		ServerProperties: input,
	}, nil
}

func getNewServer(c *core.CommandConfig) (*resources.Server, error) {
	input := resources.ServerProperties{}

	serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))

	input.SetType(serverType)
	input.SetAvailabilityZone(availabilityZone)
	input.SetName(name)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Type set: %v", serverType))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property AvailabilityZone set: %v", availabilityZone))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))

	// CUBE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverCubeType {
		input.ServerProperties.CpuFamily = nil
		if fn := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(fn) {
			// NOTE 19.07.2023:
			// In the past, all CUBE servers had to have "INTEL_SKYLAKE" as a CPU Family.
			// As such, INTEL_SKYLAKE was hardcoded as the CpuFamily field.
			//
			// However, something changed on the API side, and started throwing errs if Cpu Family was set:
			// `[VDC-5-1921] The attribute 'cpuFamily' must not be provided for Cube servers.`
			//
			// I will allow the user to modify this field, but only if the flag is explicitly set,
			// in case the API changes back to its old state in the future

			input.SetCpuFamily(viper.GetString(fn))
		}
		if !input.HasName() {
			input.SetName("Unnamed Cube")
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)) {
			templateUuid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))
			input.SetTemplateUuid(templateUuid)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property TemplateUuid set: %v", templateUuid))
		}
	}

	// ENTERPRISE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverEnterpriseType {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCpuFamily)) &&
			viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily)) != cloudapiv6.DefaultServerCPUFamily {
			input.SetCpuFamily(viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily)))
		} else {
			cpuFamily, err := DefaultCpuFamily(c)
			if err != nil {
				return nil, err
			}

			input.SetCpuFamily(cpuFamily)
		}

		if !input.HasName() {
			input.SetName("Unnamed Server")
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
			cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
			input.SetCores(cores)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Cores set: %v", cores))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
			size, err := utils2.ConvertSize(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			input.SetRam(int32(size))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB", int32(size)))
		}
	}

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagType)) == serverVCPUType {
		input.ServerProperties.CpuFamily = nil
		if fn := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(fn) {
			input.SetCpuFamily(viper.GetString(fn))
		}

		if !input.HasName() {
			input.SetName("Unnamed VCPU")
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
			cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
			input.SetCores(cores)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Cores set: %v", cores))
		}
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
			size, err := utils2.ConvertSize(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)),
				utils2.MegaBytes,
			)
			if err != nil {
				return nil, err
			}

			input.SetRam(int32(size))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Ram set: %vMB", int32(size)))
		}
	}

	return &resources.Server{
		Server: compute.Server{
			Properties: &input.ServerProperties,
		},
	}, nil
}

func getNewDAS(c *core.CommandConfig) (*resources.Volume, error) {
	volumeProper := resources.VolumeProperties{}

	volumeProper.SetType("DAS")
	volumeProper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeName)))
	volumeProper.SetBus(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus)))

	if (!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		volumeProper.SetLicenceType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) {
		volumeProper.SetImage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		volumeProper.SetImageAlias(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) {
		volumeProper.SetImagePassword(viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("SSH Key Paths: %v", sshKeysPaths))

		sshKeys, err := getSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		volumeProper.SetSshKeys(sshKeys)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SshKeys set"))
	}

	return &resources.Volume{
		Volume: compute.Volume{
			Properties: &volumeProper.VolumeProperties,
		},
	}, nil
}

func DeleteAllServers(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Servers..."))

	servers, resp, err := c.CloudApiV6Services.Servers().List(dcId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	serversItems, ok := servers.GetItemsOk()
	if !ok || serversItems == nil {
		return fmt.Errorf("could not get items of Servers")
	}

	if len(*serversItems) <= 0 {
		return fmt.Errorf("no Servers found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Servers to be deleted:"))
	for _, server := range *serversItems {
		delIdAndName := ""

		if id, ok := server.GetIdOk(); ok && id != nil {
			delIdAndName += "Server Id: " + *id
		}

		if properties, ok := server.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " Server Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Servers", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Servers..."))

	var multiErr error
	for _, server := range *serversItems {
		id, ok := server.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Starting deleting Server with id: %v from datacenter with id: %v... ", *id, dcId))

		resp, err = c.CloudApiV6Services.Servers().Delete(dcId, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Servers successfully deleted"))
	return nil
}

func DefaultCpuFamily(c *core.CommandConfig) (string, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
	if err != nil {
		return "", err
	}

	if dc.Properties == nil {
		return "", fmt.Errorf("could not retrieve Datacenter Properties")
	}

	if dc.Properties.CpuArchitecture == nil {
		return "", errors.New("could not retrieve CpuArchitecture")
	}

	cpuArch := (*dc.Properties.CpuArchitecture)[0]

	if cpuArch.CpuFamily == nil {
		return "", errors.New("could not retrieve CpuFamily")
	}

	return *cpuArch.CpuFamily, nil
}
