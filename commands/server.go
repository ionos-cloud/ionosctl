package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func server() *core.Command {
	ctx := context.TODO()
	serverCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "svr"},
			Short:            "Server Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl server` + "`" + ` allow you to create, list, get, update, delete, start, stop, reboot Servers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultServerCols,
		fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", defaultServerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(serverCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = serverCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultServerCols, cobra.ShellCompDirectiveNoFileComp
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
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
		LongDesc: `Use this command to create a Server in a specified Virtual Data Center. It is required that the number of cores for the Server and the amount of memory for the Server to be set.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways: 

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Cores
* RAM`,
		Example:    createServerExample,
		PreCmdRun:  PreRunDcIdCoresRam,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgTemplateId, "", "", "The unique Template Id for creating a CUBE Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getTemplatesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgType, "", "", "Type usages for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ENTERPRISE", "CUBE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Server")
	create.AddIntFlag(config.ArgCores, "", config.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits "+config.RequiredFlag)
	create.AddStringFlag(config.ArgRam, "", "", "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgCPUFamily, "", config.DefaultServerCPUFamily, "CPU Family for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgAvailabilityZone, config.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Server")
	update.AddStringFlag(config.ArgCPUFamily, "", config.DefaultServerCPUFamily, "CPU Family of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgAvailabilityZone, config.ArgAvailabilityZoneShort, "", "Availability zone of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(config.ArgCores, "", config.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits")
	update.AddStringFlag(config.ArgRam, "", strconv.Itoa(config.DefaultServerRAM), "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

NOTE: This will not automatically remove the storage Volume(s) attached to a Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    deleteServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server deletion to be executed")
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
	suspend.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = suspend.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = suspend.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(suspend.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	suspend.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server reboot to be executed")
	suspend.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")

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
	start.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(start.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	stop.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(stop.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	reboot.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(reboot.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	resume.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = resume.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.RequiredFlagServerId)
	_ = resume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(resume.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	resume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server reboot to be executed")
	resume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")

	serverCmd.AddCommand(serverVolume())
	serverCmd.AddCommand(serverCdrom())

	return serverCmd
}

func PreRunDcIdCoresRam(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgCores, config.ArgRam)
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgServerId)
}

func RunServerList(c *core.CommandConfig) error {
	servers, _, err := c.Servers().List(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, getServers(servers)))
}

func RunServerGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, GetStateServer, viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))); err != nil {
		return err
	}
	svr, _, err := c.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
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
	// If CPU Family has not been set, take the default value
	if !proper.ServerProperties.HasCpuFamily() {
		proper.ServerProperties.SetCpuFamily(viper.GetString(core.GetFlagName(c.NS, config.ArgCpuFamily)))
	}
	svr, resp, err := c.Servers().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		resources.Server{
			Server: ionoscloud.Server{
				Properties: &proper.ServerProperties,
			},
		},
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := svr.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, GetStateServer, *id); err != nil {
				return err
			}
			if svr, _, err = c.Servers().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
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
	input, err := getNewServerInfo(c)
	if err != nil {
		return err
	}
	svr, resp, err := c.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		*input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, GetStateServer, viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))); err != nil {
			return err
		}
		if svr, _, err = c.Servers().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getServerPrint(resp, c, []resources.Server{*svr}))
}

func RunServerDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server"); err != nil {
		return err
	}
	resp, err := c.Servers().Delete(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerStart(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "start server"); err != nil {
		return err
	}
	resp, err := c.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerStop(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "stop server"); err != nil {
		return err
	}
	resp, err := c.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerSuspend(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "suspend cube server"); err != nil {
		return err
	}
	resp, err := c.Servers().Suspend(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func RunServerReboot(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "reboot server"); err != nil {
		return err
	}
	resp, err := c.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func getNewServerInfo(c *core.CommandConfig) (*resources.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCPUFamily)) {
		input.SetCpuFamily(viper.GetString(core.GetFlagName(c.NS, config.ArgCPUFamily)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgAvailabilityZone)) {
		input.SetAvailabilityZone(viper.GetString(core.GetFlagName(c.NS, config.ArgAvailabilityZone)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCores)) {
		input.SetCores(viper.GetInt32(core.GetFlagName(c.NS, config.ArgCores)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgTemplateId)) {
		input.SetTemplateUuid(viper.GetString(core.GetFlagName(c.NS, config.ArgTemplateId)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgType)) {
		input.SetType(viper.GetString(core.GetFlagName(c.NS, config.ArgType)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgRam)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, config.ArgRam)),
			utils.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		input.SetRam(int32(size))
	}
	return &resources.ServerProperties{
		ServerProperties: input,
	}, nil
}

func RunServerResume(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "resume cube server"); err != nil {
		return err
	}
	resp, err := c.Servers().Resume(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

// Wait for State

func GetStateServer(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.Servers().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)), objId)
	if err != nil {
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

// Output Printing

var defaultServerCols = []string{"ServerId", "Name", "AvailabilityZone", "Cores", "Ram", "CpuFamily", "VmState", "State"}

type ServerPrint struct {
	ServerId         string `json:"ServerId,omitempty"`
	Name             string `json:"Name,omitempty"`
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	State            string `json:"State,omitempty"`
	Cores            int32  `json:"Cores,omitempty"`
	Ram              string `json:"Ram,omitempty"`
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
		if id, ok := s.GetIdOk(); ok && id != nil {
			serverPrint.ServerId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				serverPrint.Name = *name
			}
			if cores, ok := properties.GetCoresOk(); ok && cores != nil {
				serverPrint.Cores = *cores
			}
			if ram, ok := properties.GetRamOk(); ok && ram != nil {
				serverPrint.Ram = fmt.Sprintf("%vMB", *ram)
			}
			if cpuFamily, ok := properties.GetCpuFamilyOk(); ok && cpuFamily != nil {
				serverPrint.CpuFamily = *cpuFamily
			}
			if zone, ok := properties.GetAvailabilityZoneOk(); ok && zone != nil {
				serverPrint.AvailabilityZone = *zone
			}
			if vmState, ok := properties.GetVmStateOk(); ok && vmState != nil {
				serverPrint.VmState = *vmState
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				serverPrint.State = *state
			}
		}
		o := structs.Map(serverPrint)
		out = append(out, o)
	}
	return out
}

func getServersIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(clientSvc.Get(), context.TODO())
	servers, _, err := serverSvc.List(datacenterId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := servers.Servers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
