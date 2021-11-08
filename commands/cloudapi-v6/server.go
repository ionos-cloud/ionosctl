package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	serverCubeType       = "CUBE"
	serverEnterpriseType = "ENTERPRISE"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultServerCols, printer.ColsMessage(allServerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(serverCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = serverCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allServerCols, cobra.ShellCompDirectiveNoFileComp
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
		LongDesc:   "Use this command to list Servers from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listServerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunServerList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, fmt.Sprintf("Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1:VALUE1,KEY2:VALUE2. Available filters: %v", completer.DataCentersFilters()))
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersFilters(), cobra.ShellCompDirectiveNoFileComp
	})

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
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Server to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for waiting for Server to be in AVAILABLE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, serverCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "server",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Server",
		LongDesc: `Use this command to create an ENTERPRISE or CUBE Server in a specified Virtual Data Center. 

* For ENTERPRISE Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

To see which CPU Family are available in which location, use ` + "`" + `ionosctl location` + "`" + ` commands.

Required values to create a Server of type ENTERPRISE:

* Data Center Id
* Cores
* RAM

* For CUBE Servers:

Servers of type CUBE will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use ` + "`" + `ionosctl template` + "`" + ` commands.

Required values to create a Server of type CUBE:

* Data Center Id
* Type
* Template Id

By default, Licence Type for Direct Attached Storage is set to LINUX. You can set it using the ` + "`" + `--licence-type` + "`" + ` option or set an Image Id. For Image Id, it is needed to set a password or SSH keys.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.`,
		Example:    createServerExample,
		PreCmdRun:  PreRunServerCreate,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Server", "Name of the Server")
	create.AddIntFlag(cloudapiv6.ArgCores, "", cloudapiv6.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgCPUFamily, "", cloudapiv6.DefaultServerCPUFamily, "CPU Family for the Server. For CUBE Servers, the CPU Family is INTEL_SKYLAKE")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgAvailabilityZone, cloudapiv6.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgTemplateId, "", "", "[CUBE Server] The unique Template Id", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TemplatesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgType, "", serverEnterpriseType, "Type usages for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{serverEnterpriseType, serverCubeType}, cobra.ShellCompDirectiveNoFileComp
	})

	// Volume Properties - for DAS Volume associated with Cube Server
	create.AddStringFlag(cloudapiv6.ArgVolumeName, "N", "Unnamed Direct Attached Storage", "[CUBE Server] Name of the Direct Attached Storage")
	create.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "[CUBE Server] The bus type of the Direct Attached Storage")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgLicenceType, "l", "LINUX", "[CUBE Server] Licence Type of the Direct Attached Storage")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"LINUX", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "[CUBE Server] The Image Alias to use instead of Image Id for the Direct Attached Storage")
	create.AddStringFlag(cloudapiv6.ArgImageId, "", "", "[CUBE Server] The Image Id or snapshot Id to be used as for the Direct Attached Storage")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "[CUBE Server] Initial image password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	create.AddStringSliceFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, []string{""}, "[CUBE Server] Absolute paths for the SSH Keys of the Direct Attached Storage")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server creation to be executed")
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for new Server to be in AVAILABLE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds]")

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
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgVolumeId, "", "", "The unique Volume Id for the BootVolume. The Volume needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgCdromId, "", "", "The unique Cdrom Id for the BootCdrom. The Cdrom needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesIdsCustom(os.Stderr, resources.ListQueryParams{Filters: &map[string]string{
			"type": "CDROM",
		}}), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Server")
	update.AddStringFlag(cloudapiv6.ArgCPUFamily, "", cloudapiv6.DefaultServerCPUFamily, "CPU Family of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgAvailabilityZone, cloudapiv6.ArgAvailabilityZoneShort, "", "Availability zone of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.ArgCores, "", cloudapiv6.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits")
	update.AddStringFlag(cloudapiv6.ArgRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server update to be executed")
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the updated Server to be in AVAILABLE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server update/for Server to be in AVAILABLE state [seconds]")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Servers form a virtual Datacenter.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server deletion [seconds]")

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
	suspend.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = suspend.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(os.Stderr, viper.GetString(core.GetFlagName(suspend.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string]string{
					"type": serverCubeType,
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server suspend to be executed")
	suspend.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server suspend [seconds]")

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
	start.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(start.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server start to be executed")
	start.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server start [seconds]")

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
	stop.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(stop.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server stop to be executed")
	stop.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server stop [seconds]")

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
	reboot.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(reboot.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server reboot to be executed")
	reboot.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")

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
	resume.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddStringFlag(cloudapiv6.ArgServerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = resume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIdsCustom(os.Stderr, viper.GetString(core.GetFlagName(resume.NS, cloudapiv6.ArgDataCenterId)),
			resources.ListQueryParams{
				Filters: &map[string]string{
					"type": serverCubeType,
				},
			}), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server resume to be executed")
	resume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server resume [seconds]")

	serverCmd.AddCommand(ServerTokenCmd())
	serverCmd.AddCommand(ServerConsoleCmd())
	serverCmd.AddCommand(ServerVolumeCmd())
	serverCmd.AddCommand(ServerCdromCmd())

	return serverCmd
}

func PreRunServerCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgCores, cloudapiv6.ArgRam},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgType, cloudapiv6.ArgTemplateId})
	if err != nil {
		return err
	}
	// Validate flags
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		err = core.CheckRequiredFlagsSets(c.Command, c.NS,
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgCores, cloudapiv6.ArgRam, cloudapiv6.ArgImageId, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgCores, cloudapiv6.ArgRam, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgCores, cloudapiv6.ArgRam, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgCores, cloudapiv6.ArgRam, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgType, cloudapiv6.ArgTemplateId, cloudapiv6.ArgImageId, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgType, cloudapiv6.ArgTemplateId, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgType, cloudapiv6.ArgTemplateId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgType, cloudapiv6.ArgTemplateId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths},
		)
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

func RunServerList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	servers, resp, err := c.CloudApiV6Services.Servers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, getServers(servers)))
}

func RunServerGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Server with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)))
	if err := utils.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
		return err
	}
	svr, resp, err := c.CloudApiV6Services.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, []resources.Server{*svr}))
}

func RunServerCreate(c *core.CommandConfig) error {
	input, err := getNewServer(c)
	if err != nil {
		return err
	}
	// If Server is of type CUBE, it will create an attached Volume
	if viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)) == serverCubeType {
		// Volume Properties
		volumeDAS, err := getNewDAS(c)
		if err != nil {
			return err
		}
		// Attach Storage
		input.SetEntities(ionoscloud.ServerEntities{
			Volumes: &ionoscloud.AttachedVolumes{
				Items: &[]ionoscloud.Volume{volumeDAS.Volume},
			},
		})
	}
	svr, resp, err := c.CloudApiV6Services.Servers().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := svr.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.ServerStateInterrogator, *id); err != nil {
				return err
			}
			if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
				*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new server id")
		}
	}
	return c.Printer.Print(getServerPrint(resp, c, []resources.Server{*svr}))
}

func RunServerUpdate(c *core.CommandConfig) error {
	input, err := getUpdateServerInfo(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating Server with ID: %v in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	svr, resp, err := c.CloudApiV6Services.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		*input,
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
			return err
		}
		if svr, _, err = c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getServerPrint(resp, c, []resources.Server{*svr}))
}

func RunServerDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll))
	if allFlag {
		resp, err = DeleteAllServers(c)
		if err != nil {
			return err
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Server with id: %v from datacenter with id: %v... ", serverId, dcId)
		resp, err = c.CloudApiV6Services.Servers().Delete(dcId, serverId)
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerStart(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "start server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server is starting... ")
	resp, err := c.CloudApiV6Services.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerStop(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "stop server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server is stopping... ")
	resp, err := c.CloudApiV6Services.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerSuspend(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "suspend cube server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server is Suspending... ")
	resp, err := c.CloudApiV6Services.Servers().Suspend(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerReboot(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "reboot server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server is rebooting... ")
	resp, err := c.CloudApiV6Services.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerResume(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "resume cube server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server is resuming... ")
	resp, err := c.CloudApiV6Services.Servers().Resume(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func getUpdateServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property name set: %v ", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCPUFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCPUFamily))
		c.Printer.Verbose("Property CpuFamily set: %v ", cpuFamily)
		input.SetCpuFamily(cpuFamily)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAvailabilityZone))
		c.Printer.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
		input.SetAvailabilityZone(availabilityZone)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCores))
		c.Printer.Verbose("Property Cores set: %v ", cores)
		input.SetCores(cores)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)) {
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
		c.Printer.Verbose("Property BootVolume set: %v ", volumeId)
		input.SetBootVolume(ionoscloud.ResourceReference{
			Id: &volumeId,
		})
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)) {
		cdromId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId))
		c.Printer.Verbose("Property BootCdrom set: %v ", cdromId)
		input.SetBootCdrom(ionoscloud.ResourceReference{
			Id: &cdromId,
		})
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRam)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRam)),
			utils.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		c.Printer.Verbose("Property Ram set: %vMB ", int32(size))
		input.SetRam(int32(size))
	}
	return &resources.ServerProperties{
		ServerProperties: input,
	}, nil
}

func getNewServer(c *core.CommandConfig) (*resources.Server, error) {
	input := resources.ServerProperties{}
	serverType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAvailabilityZone))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	input.SetType(serverType)
	input.SetAvailabilityZone(availabilityZone)
	input.SetName(name)
	c.Printer.Verbose("Property Type set: %v", serverType)
	c.Printer.Verbose("Property AvailabilityZone set: %v", availabilityZone)
	c.Printer.Verbose("Property Name set: %v", name)

	// CUBE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)) == serverCubeType {
		// Right now, for the CUBE Server - only INTEL_SKYLAKE is supported
		input.SetCpuFamily("INTEL_SKYLAKE")
		if !input.HasName() {
			input.SetName("Unnamed Cube")
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId)) {
			templateUuid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTemplateId))
			input.SetTemplateUuid(templateUuid)
			c.Printer.Verbose("Property TemplateUuid set: %v", templateUuid)
		}
	}

	// ENTERPRISE Server Properties
	if viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)) == serverEnterpriseType {
		input.SetCpuFamily(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCpuFamily)))
		if !input.HasName() {
			input.SetName("Unnamed Server")
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCores)) {
			cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCores))
			input.SetCores(cores)
			c.Printer.Verbose("Property Cores set: %v", cores)
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRam)) {
			size, err := utils.ConvertSize(
				viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRam)),
				utils.MegaBytes,
			)
			if err != nil {
				return nil, err
			}
			input.SetRam(int32(size))
			c.Printer.Verbose("Property Ram set: %vMB", int32(size))
		}
	}
	return &resources.Server{
		Server: ionoscloud.Server{
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
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPassword)) {
		volumeProper.SetImagePassword(viper.GetString(core.GetFlagName(c.NS, config.ArgPassword)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))
		c.Printer.Verbose("SSH Key Paths: %v", sshKeysPaths)
		sshKeys, err := getSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}
		volumeProper.SetSshKeys(sshKeys)
		c.Printer.Verbose("Property SshKeys set")
	}
	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &volumeProper.VolumeProperties,
		},
	}, nil
}

func DeleteAllServers(c *core.CommandConfig) (*resources.Response, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	_ = c.Printer.Print("Servers to be deleted:")
	servers, resp, err := c.CloudApiV6Services.Servers().List(dcId, resources.ListQueryParams{})
	if err != nil {
		return nil, err
	}
	if serversItems, ok := servers.GetItemsOk(); ok && serversItems != nil {
		for _, server := range *serversItems {
			if id, ok := server.GetIdOk(); ok && id != nil {
				_ = c.Printer.Print("Server Id: " + *id)
			}
			if properties, ok := server.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print("Server Name: " + *name)
				}
			}
		}
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Servers"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Servers...")

		for _, server := range *serversItems {
			if id, ok := server.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting Server with id: %v from datacenter with id: %v... ", *id, dcId)
				resp, err = c.CloudApiV6Services.Servers().Delete(dcId, *id)
				if resp != nil {
					c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
				}
				if err != nil {
					return nil, err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, err
}

// Output Printing

var (
	defaultServerCols = []string{"ServerId", "Name", "Type", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State"}
	allServerCols     = []string{"ServerId", "Name", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State", "TemplateId", "Type", "BootCdromId", "BootVolumeId"}
)

type ServerPrint struct {
	ServerId         string `json:"ServerId,omitempty"`
	Name             string `json:"Name,omitempty"`
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	State            string `json:"State,omitempty"`
	Cores            int32  `json:"Cores,omitempty"`
	Ram              string `json:"Ram,omitempty"`
	CpuFamily        string `json:"CpuFamily,omitempty"`
	VmState          string `json:"VmState,omitempty"`
	BootVolumeId     string `json:"BootVolumeId,omitempty"`
	BootCdromId      string `json:"BootCdromId,omitempty"`
	TemplateId       string `json:"TemplateId,omitempty"`
	Type             string `json:"Type,omitempty"`
}

func getServerPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.Server) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getServersKVMaps(ss)
			r.Columns = getServersCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getServersCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultServerCols
	}

	columnsMap := map[string]string{
		"ServerId":         "ServerId",
		"Name":             "Name",
		"AvailabilityZone": "AvailabilityZone",
		"State":            "State",
		"VmState":          "VmState",
		"Cores":            "Cores",
		"Ram":              "Ram",
		"CpuFamily":        "CpuFamily",
		"TemplateId":       "TemplateId",
		"Type":             "Type",
		"BootVolumeId":     "BootVolumeId",
		"BootCdromId":      "BootCdromId",
	}
	var serverCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			serverCols = append(serverCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return serverCols
}

func getServers(servers resources.Servers) []resources.Server {
	ss := make([]resources.Server, 0)
	for _, s := range *servers.Items {
		ss = append(ss, resources.Server{Server: s})
	}
	return ss
}

func getServersKVMaps(ss []resources.Server) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var serverPrint ServerPrint
		if idOk, ok := s.GetIdOk(); ok && idOk != nil {
			serverPrint.ServerId = *idOk
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if nameOk, ok := properties.GetNameOk(); ok && nameOk != nil {
				serverPrint.Name = *nameOk
			}
			if coresOk, ok := properties.GetCoresOk(); ok && coresOk != nil {
				serverPrint.Cores = *coresOk
			}
			if ramOk, ok := properties.GetRamOk(); ok && ramOk != nil {
				serverPrint.Ram = fmt.Sprintf("%vMB", *ramOk)
			}
			if cpuFamilyOk, ok := properties.GetCpuFamilyOk(); ok && cpuFamilyOk != nil {
				serverPrint.CpuFamily = *cpuFamilyOk
			}
			if zoneOk, ok := properties.GetAvailabilityZoneOk(); ok && zoneOk != nil {
				serverPrint.AvailabilityZone = *zoneOk
			}
			if vmStateOk, ok := properties.GetVmStateOk(); ok && vmStateOk != nil {
				serverPrint.VmState = *vmStateOk
			}
			if templateUuidOk, ok := properties.GetTemplateUuidOk(); ok && templateUuidOk != nil {
				serverPrint.TemplateId = *templateUuidOk
			}
			if typeOk, ok := properties.GetTypeOk(); ok && typeOk != nil {
				serverPrint.Type = *typeOk
			}
			if bootVolumeOk, ok := properties.GetBootVolumeOk(); ok && bootVolumeOk != nil {
				if idOk, ok := bootVolumeOk.GetIdOk(); ok && idOk != nil {
					serverPrint.BootVolumeId = *idOk
				}
			}
			if bootCdromOk, ok := properties.GetBootCdromOk(); ok && bootCdromOk != nil {
				if idOk, ok := bootCdromOk.GetIdOk(); ok && idOk != nil {
					serverPrint.BootCdromId = *idOk
				}
			}
		}
		if metadataOk, ok := s.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				serverPrint.State = *stateOk
			}
		}
		o := structs.Map(serverPrint)
		out = append(out, o)
	}
	return out
}
