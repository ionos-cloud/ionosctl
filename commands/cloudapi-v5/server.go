package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerCmd() *core.Command {
	ctx := context.TODO()
	serverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "svr"},
			Short:            "Server Operations",
			Long:             "The sub-commands of `ionosctl server` allow you to create, list, get, update, delete, start, stop, reboot Servers.",
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
		LongDesc:   "Use this command to list Servers from a specified Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ServersFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listServerExample,
		PreCmdRun:  PreRunServerList,
		CmdRun:     RunServerList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	get.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
		LongDesc: `Use this command to create a Server in a specified Virtual Data Center. It is required that the number of cores and the amount of memory for the Server to be set.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Cores
* RAM`,
		Example:    createServerExample,
		PreCmdRun:  PreRunDcIdServerCoresRam,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "Unnamed Server", "Name of the Server")
	create.AddIntFlag(cloudapiv5.ArgCores, "", cloudapiv5.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv5.ArgRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgCpuFamily, "", cloudapiv5.DefaultServerCPUFamily, "CPU Family for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgAvailabilityZone, cloudapiv5.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
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

Note: The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    updateServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "Name of the Server")
	update.AddStringFlag(cloudapiv5.ArgCpuFamily, "", "", "CPU Family of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgCpuFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgAvailabilityZone, cloudapiv5.ArgAvailabilityZoneShort, "", "Availability zone of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv5.ArgCores, "", cloudapiv5.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits")
	update.AddStringFlag(cloudapiv5.ArgRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgVolumeId, "", "", "The unique Volume Id for the BootVolume. The Volume needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgCdromId, "", "", "The unique Cdrom Id for the BootCdrom. The Cdrom needs to be already attached to the Server")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesIdsCustom(os.Stderr, resources.ListQueryParams{Filters: &map[string]string{
			"type": "CDROM",
		}}), cobra.ShellCompDirectiveNoFileComp
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
	deleteCmd.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all Servers form a virtual Datacenter.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server deletion [seconds]")

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
	start.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(start.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	stop.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(stop.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
		Example:    resetServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerReboot,
		InitClient: true,
	})
	reboot.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddStringFlag(cloudapiv5.ArgServerId, cloudapiv5.ArgIdShort, "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(reboot.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server reboot to be executed")
	reboot.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")

	serverCmd.AddCommand(ServerVolumeCmd())
	serverCmd.AddCommand(ServerCdromCmd())

	return serverCmd
}

func PreRunServerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.ServersFilters(), completer.ServersFiltersUsage())
	}
	return nil
}

func PreRunDcIdServerCoresRam(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgCores, cloudapiv5.ArgRam)
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId)
}

func PreRunDcServerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId},
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgAll},
	)
}

func RunServerList(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters
			if val, ok := filters["ram"]; ok {
				convertedSize, err := utils.ConvertSize(val, utils.MegaBytes)
				if err != nil {
					return err
				}
				filters["ram"] = strconv.Itoa(convertedSize)
				listQueryParams.Filters = &filters
			}
		}
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	servers, resp, err := c.CloudApiV5Services.Servers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, getServers(servers)))
}

func RunServerGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	c.Printer.Verbose("Server with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)))
	if err := utils.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))); err != nil {
		return err
	}
	svr, resp, err := c.CloudApiV5Services.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, []resources.Server{*svr}))
}

func RunServerCreate(c *core.CommandConfig) error {
	proper, err := getNewServerInfo(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Creating Server in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	svr, resp, err := c.CloudApiV5Services.Servers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		resources.Server{
			Server: ionoscloud.Server{
				Properties: &proper.ServerProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
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
			if svr, _, err = c.CloudApiV5Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
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
	input, err := getServerInfo(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating Server with ID: %v in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	svr, resp, err := c.CloudApiV5Services.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
		*input,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, waiter.ServerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))); err != nil {
			return err
		}
		if svr, _, err = c.CloudApiV5Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getServerPrint(resp, c, []resources.Server{*svr}))
}

func RunServerDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll)) {
		if err := DeleteAllServers(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Server with id: %v from datacenter with id: %v... ", serverId, dcId)
		resp, err := c.CloudApiV5Services.Servers().Delete(dcId, serverId)
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getServerPrint(resp, c, nil))
	}
}

func RunServerStart(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "start server"); err != nil {
		return err
	}
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is starting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)))
	resp, err := c.CloudApiV5Services.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
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
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is stopping... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)))
	resp, err := c.CloudApiV5Services.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
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
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is rebooting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)))
	resp, err := c.CloudApiV5Services.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func getNewServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}

	// Setting Properties for the New Server
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
	c.Printer.Verbose("Property name set: %v ", name)
	input.SetName(name)
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgCpuFamily))
	c.Printer.Verbose("Property CpuFamily set: %v ", cpuFamily)
	input.SetCpuFamily(cpuFamily)
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgAvailabilityZone))
	c.Printer.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
	input.SetAvailabilityZone(availabilityZone)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgCores))
		c.Printer.Verbose("Property Cores set: %v ", cores)
		input.SetCores(cores)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgRam)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgRam)),
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

func getServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
		c.Printer.Verbose("Property name set: %v ", name)
		input.SetName(name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgCpuFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgCpuFamily))
		c.Printer.Verbose("Property CpuFamily set: %v ", cpuFamily)
		input.SetCpuFamily(cpuFamily)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgAvailabilityZone))
		c.Printer.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
		input.SetAvailabilityZone(availabilityZone)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgCores))
		c.Printer.Verbose("Property Cores set: %v ", cores)
		input.SetCores(cores)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId)) {
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId))
		c.Printer.Verbose("Property BootVolume set: %v ", volumeId)
		input.SetBootVolume(ionoscloud.ResourceReference{
			Id: &volumeId,
		})
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgCdromId)) {
		cdromId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgCdromId))
		c.Printer.Verbose("Property BootCdrom set: %v ", cdromId)
		input.SetBootCdrom(ionoscloud.ResourceReference{
			Id: &cdromId,
		})
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgRam)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgRam)),
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

func DeleteAllServers(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting Servers...")
	servers, _, err := c.CloudApiV5Services.Servers().List(dcId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if serversItems, ok := servers.GetItemsOk(); ok && serversItems != nil {
		if len(*serversItems) > 0 {
			_ = c.Printer.Print("Servers to be deleted:")
			for _, server := range *serversItems {
				toPrint := ""
				if id, ok := server.GetIdOk(); ok && id != nil {
					toPrint += "Server Id: " + *id
				}
				if properties, ok := server.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += " Server Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Servers"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Servers...")
			var multiErr error
			for _, server := range *serversItems {
				if id, ok := server.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Server with id: %v from datacenter with id: %v... ", *id, dcId)
					resp, err := c.CloudApiV5Services.Servers().Delete(dcId, *id)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return err
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Servers found")
		}
	} else {
		return errors.New("could not get items of Servers")
	}
}

// Output Printing

var (
	defaultServerCols = []string{"ServerId", "Name", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State"}
	allServerCols     = []string{"ServerId", "Name", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State", "BootVolumeId", "BootCdromId"}
)

type ServerPrint struct {
	ServerId         string `json:"ServerId,omitempty"`
	Name             string `json:"Name,omitempty"`
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	State            string `json:"State,omitempty"`
	Cores            int32  `json:"Cores,omitempty"`
	Ram              string `json:"Ram,omitempty"`
	BootVolumeId     string `json:"BootVolumeId,omitempty"`
	BootCdromId      string `json:"BootCdromId,omitempty"`
	CpuFamily        string `json:"CpuFamily,omitempty"`
	VmState          string `json:"VmState,omitempty"`
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
	if items, ok := servers.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.Server{Server: s})
		}
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
		if propertiesOk, ok := s.GetPropertiesOk(); ok && propertiesOk != nil {
			if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
				serverPrint.Name = *nameOk
			}
			if coresOk, ok := propertiesOk.GetCoresOk(); ok && coresOk != nil {
				serverPrint.Cores = *coresOk
			}
			if ramOk, ok := propertiesOk.GetRamOk(); ok && ramOk != nil {
				serverPrint.Ram = fmt.Sprintf("%vMB", *ramOk)
			}
			if cpuFamilyOk, ok := propertiesOk.GetCpuFamilyOk(); ok && cpuFamilyOk != nil {
				serverPrint.CpuFamily = *cpuFamilyOk
			}
			if zoneOk, ok := propertiesOk.GetAvailabilityZoneOk(); ok && zoneOk != nil {
				serverPrint.AvailabilityZone = *zoneOk
			}
			if vmStateOk, ok := propertiesOk.GetVmStateOk(); ok && vmStateOk != nil {
				serverPrint.VmState = *vmStateOk
			}
			if bootVolumeOk, ok := propertiesOk.GetBootVolumeOk(); ok && bootVolumeOk != nil {
				if idOk, ok := bootVolumeOk.GetIdOk(); ok && idOk != nil {
					serverPrint.BootVolumeId = *idOk
				}
			}
			if bootCdromOk, ok := propertiesOk.GetBootCdromOk(); ok && bootCdromOk != nil {
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
