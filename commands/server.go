package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func server() *builder.Command {
	serverCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "server",
			Short:            "Server Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl server` + "`" + ` allow you to create, list, get, update, delete, start, stop, reboot Servers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(serverCmd.Command.Use, config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = serverCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdValidate, RunServerList, "list", "List Servers",
		"Use this command to list Servers from a specified Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		listServerExample, true)

	get := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerGet, "get", "Get a Server",
		"Use this command to get information about a specified Server from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		getServerExample, true)
	get.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	create := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdValidate, RunServerCreate, "create", "Create a Server",
		`Use this command to create a Server in a specified Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, createServerExample, true)
	create.AddStringFlag(config.ArgServerName, "", "", "Name of the Server")
	create.AddIntFlag(config.ArgServerCores, "", config.DefaultServerCores, "Cores option of the Server")
	create.AddIntFlag(config.ArgServerRAM, "", config.DefaultServerRAM, "RAM[GB] option for the Server")
	create.AddStringFlag(config.ArgServerCPUFamily, "", config.DefaultServerCPUFamily, "CPU Family for the Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgServerCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgServerZone, "", "AUTO", "Availability zone of the Server")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgServerZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be created [seconds]")

	update := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerUpdate, "update", "Update a Server",
		`Use this command to update a specified Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, updateServerExample, true)
	update.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgServerName, "", "", "Name of the Server")
	update.AddStringFlag(config.ArgServerCPUFamily, "", config.DefaultServerCPUFamily, "CPU Family of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgServerCPUFamily, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgServerZone, "", "", "Availability zone of the Server")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgServerZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(config.ArgServerCores, "", config.DefaultServerCores, "Cores option of the Server")
	update.AddIntFlag(config.ArgServerRAM, "", config.DefaultServerRAM, "RAM[GB] option for the Server")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be updated [seconds]")

	deleteCmd := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerDelete, "delete", "Delete a Server",
		`Use this command to delete a specified Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, deleteServerExample, true)
	deleteCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be deleted [seconds]")

	start := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerStart, "start", "Start a Server",
		`Use this command to start specified Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, startServerExample, true)
	start.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to start")
	start.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be started [seconds]")

	stop := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerStop, "stop", "Stop a Server",
		`Use this command to stop specified Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, stopServerExample, true)
	stop.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Server to stop")
	stop.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be stopped [seconds]")

	reboot := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerReboot, "reboot", "Force a hard reboot of a Server",
		`Use this command to force a hard reboot of the Server. Do not use this method if you want to gracefully reboot the machine. This is the equivalent of powering off the machine and turning it back on.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, resetServerExample, true)
	reboot.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	attachVolume := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerVolumeIdsValidate, RunServerAttachVolume, "attach-volume", "Attach a Volume to a Server",
		`Use this command to attach a Volume to a Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`, attachVolumeServerExample, true)
	attachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to attach to Server")
	attachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be attached to a Server [seconds]")

	listAttached := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerListVolumes, "list-volumes", "List attached Volumes from a Server",
		"Use this command to get a list of attached Volumes to a Server from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		listVolumesServerExample, true)
	listAttached.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = listAttached.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	getAttached := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerVolumeIdsValidate, RunServerGetVolume, "get-volume", "Get an attached Volume from a Server",
		"Use this command to retrieve information about an attached Volume on Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Volume Id",
		getVolumeServerExample, true)
	getAttached.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getAttached.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(serverCmd.Name(), getAttached.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	detachVolume := builder.NewCommand(context.TODO(), serverCmd, PreRunGlobalDcIdServerVolumeIdsValidate, RunServerDetachVolume, "detach-volume", "Detach a Volume from a Server",
		`Use this command to detach a Volume from a Server.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`, detachVolumeServerExample, true)
	detachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(serverCmd.Name(), detachVolume.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to detach from Server")
	detachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be detached from a Server [seconds]")

	labelServer(serverCmd)

	return serverCmd
}

func PreRunGlobalDcIdServerIdValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalDcIdServerVolumeIdsValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId, config.ArgVolumeId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunServerList(c *builder.CommandConfig) error {
	servers, _, err := c.Servers().List(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	ss := getServers(servers)
	return c.Printer.Print(printer.Result{
		OutputJSON: servers,
		KeyValue:   getServersKVMaps(ss),
		Columns:    getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunServerGet(c *builder.CommandConfig) error {
	server, _, err := c.Servers().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: server,
		KeyValue:   getServersKVMaps([]resources.Server{*server}),
		Columns:    getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunServerCreate(c *builder.CommandConfig) error {
	server, resp, err := c.Servers().Create(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerName)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCPUFamily)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerZone)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCores)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerRAM)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  server,
		KeyValue:    getServersKVMaps([]resources.Server{*server}),
		Columns:     getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerUpdate(c *builder.CommandConfig) error {
	input := resources.ServerProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCPUFamily)) {
		input.SetCpuFamily(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCPUFamily)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerZone)) {
		input.SetAvailabilityZone(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerZone)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCores)) {
		input.SetCores(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCores)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerRAM)) {
		input.SetRam(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerRAM)))
	}
	server, resp, err := c.Servers().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		input,
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		KeyValue:    getServersKVMaps([]resources.Server{*server}),
		Columns:     getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		OutputJSON:  server,
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "update",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server")
	if err != nil {
		return err
	}
	resp, err := c.Servers().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerStart(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "start server")
	if err != nil {
		return err
	}
	resp, err := c.Servers().Start(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "start",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerStop(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "stop server")
	if err != nil {
		return err
	}
	resp, err := c.Servers().Stop(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "stop",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerReboot(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "reboot server")
	if err != nil {
		return err
	}
	resp, err := c.Servers().Reboot(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "server",
		Verb:        "reboot",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunServerAttachVolume(c *builder.CommandConfig) error {
	attachedVol, resp, err := c.Servers().AttachVolume(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(attachedVol)))
}

func RunServerListVolumes(c *builder.CommandConfig) error {
	attachedVols, _, err := c.Servers().ListVolumes(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getAttachedVolumes(attachedVols)))
}

func RunServerGetVolume(c *builder.CommandConfig) error {
	attachedVolume, _, err := c.Servers().GetVolume(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(attachedVolume)))
}

func RunServerDetachVolume(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume from server")
	if err != nil {
		return err
	}
	resp, err := c.Servers().DetachVolume(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

// Output Printing

var defaultServerCols = []string{"ServerId", "Name", "AvailabilityZone", "State", "Cores", "Ram", "CpuFamily"}

type ServerPrint struct {
	ServerId         string `json:"ServerId,omitempty"`
	Name             string `json:"Name,omitempty"`
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	State            string `json:"State,omitempty"`
	Cores            int32  `json:"Cores,omitempty"`
	Ram              string `json:"Ram,omitempty"`
	CpuFamily        string `json:"CpuFamily,omitempty"`
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
		properties := s.GetProperties()
		metadata := s.GetMetadata()
		var serverPrint ServerPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			serverPrint.ServerId = *id
		}
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
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			serverPrint.State = *state
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
