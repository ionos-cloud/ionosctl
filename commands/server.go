package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
			Long:             "The sub-commands of `ionosctl server` allow you to create, list, get, update, delete, start, stop, reboot Servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultServerCols, utils.ColsMessage(defaultServerCols))
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
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
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
		LongDesc: `Use this command to create a Server in a specified Virtual Data Center. You can set the number of cores and the amount of memory for the Server.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.` + "`" + `--ram 256` + "`" + ` equals 256MB.
* providing both the value and the unit, e.g.` + "`" + `--ram 1GB` + "`" + `.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can also wait for Server to be in AVAILABLE state using ` + "`" + `--wait-for-state` + "`" + ` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id`,
		Example:    createServerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunServerCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Unnamed Server", "Name of the Server")
	create.AddIntFlag(config.ArgCores, "", config.DefaultServerCores, "The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits")
	create.AddStringFlag(config.ArgRam, "", config.DefaultServerRAM, "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
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
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
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
	update.AddStringFlag(config.ArgRam, "", config.DefaultServerRAM, "The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB")
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

NOTE: This will not automatically remove the storage Volumes attached to a Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id`,
		Example:    deleteServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server deletion to be executed")
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
	start.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
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
	stop.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
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
		Example:    resetServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerReboot,
		InitClient: true,
	})
	reboot.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddStringFlag(config.ArgServerId, config.ArgIdShort, "", config.ServerId, core.RequiredFlagOption())
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(reboot.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Server reboot to be executed")
	reboot.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Server reboot [seconds]")

	serverCmd.AddCommand(serverVolume())
	serverCmd.AddCommand(serverCdrom())

	return serverCmd
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgDataCenterId, config.ArgServerId)
}

func RunServerList(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	servers, resp, err := c.Servers().List(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, getServers(servers)))
}

func RunServerGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	c.Printer.Verbose("Server with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)))
	if err := utils.WaitForState(c, GetStateServer, viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))); err != nil {
		return err
	}
	svr, resp, err := c.Servers().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(nil, c, []v5.Server{*svr}))
}

func RunServerCreate(c *core.CommandConfig) error {
	proper, err := getNewServerInfo(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Creating Server in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	svr, resp, err := c.Servers().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		v5.Server{
			Server: ionoscloud.Server{
				Properties: &proper.ServerProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	return c.Printer.Print(getServerPrint(resp, c, []v5.Server{*svr}))
}

func RunServerUpdate(c *core.CommandConfig) error {
	input, err := getServerInfo(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating Server with ID: %v in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	svr, resp, err := c.Servers().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		*input,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	return c.Printer.Print(getServerPrint(resp, c, []v5.Server{*svr}))
}

func RunServerDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server"); err != nil {
		return err
	}
	c.Printer.Verbose("Server with id: %v from datacenter with id: %v is deleting... ", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)), viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	resp, err := c.Servers().Delete(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is starting... ", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)))
	resp, err := c.Servers().Start(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is stopping... ", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)))
	resp, err := c.Servers().Stop(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	c.Printer.Verbose("Server with ID: %v is rebooting... ", viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)))
	resp, err := c.Servers().Reboot(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getServerPrint(resp, c, nil))
}

func getNewServerInfo(c *core.CommandConfig) (*v5.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}

	// Setting Properties for the New Server
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	c.Printer.Verbose("Property name set: %v ", name)
	input.SetName(name)
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, config.ArgCPUFamily))
	c.Printer.Verbose("Property CpuFamily set: %v ", cpuFamily)
	input.SetCpuFamily(cpuFamily)
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, config.ArgAvailabilityZone))
	c.Printer.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
	input.SetAvailabilityZone(availabilityZone)
	cores := viper.GetInt32(core.GetFlagName(c.NS, config.ArgCores))
	c.Printer.Verbose("Property Cores set: %v ", cores)
	input.SetCores(cores)
	size, err := utils.ConvertSize(
		viper.GetString(core.GetFlagName(c.NS, config.ArgRam)),
		utils.MegaBytes,
	)
	if err != nil {
		return nil, err
	}
	c.Printer.Verbose("Property Ram set: %vMB ", int32(size))
	input.SetRam(int32(size))

	return &v5.ServerProperties{
		ServerProperties: input,
	}, nil
}

func getServerInfo(c *core.CommandConfig) (*v5.ServerProperties, error) {
	input := ionoscloud.ServerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
		c.Printer.Verbose("Property name set: %v ", name)
		input.SetName(name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCPUFamily)) {
		cpuFamily := viper.GetString(core.GetFlagName(c.NS, config.ArgCPUFamily))
		c.Printer.Verbose("Property CpuFamily set: %v ", cpuFamily)
		input.SetCpuFamily(cpuFamily)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgAvailabilityZone)) {
		availabilityZone := viper.GetString(core.GetFlagName(c.NS, config.ArgAvailabilityZone))
		c.Printer.Verbose("Property AvailabilityZone set: %v ", availabilityZone)
		input.SetAvailabilityZone(availabilityZone)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, config.ArgCores))
		c.Printer.Verbose("Property Cores set: %v ", cores)
		input.SetCores(cores)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgRam)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, config.ArgRam)),
			utils.MegaBytes,
		)
		if err != nil {
			return nil, err
		}
		c.Printer.Verbose("Property Ram set: %vMB ", int32(size))
		input.SetRam(int32(size))
	}
	return &v5.ServerProperties{
		ServerProperties: input,
	}, nil
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

func getServerPrint(resp *v5.Response, c *core.CommandConfig, ss []v5.Server) printer.Result {
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

func getServers(servers v5.Servers) []v5.Server {
	ss := make([]v5.Server, 0)
	for _, s := range *servers.Items {
		ss = append(ss, v5.Server{Server: s})
	}
	return ss
}

func getServersKVMaps(ss []v5.Server) []map[string]interface{} {
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
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	serverSvc := v5.NewServerService(clientSvc.Get(), context.TODO())
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
