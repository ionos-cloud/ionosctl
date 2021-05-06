package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
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
	ctx := context.TODO()
	serverCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "server",
			Short:            "Server Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl server` + "`" + ` allow you to create, list, get, update, delete, start, stop, reboot Servers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := serverCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultServerCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(serverCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	list := builder.NewCommand(ctx, serverCmd, PreRunDataCenterId, RunServerList, "list", "List Servers",
		"Use this command to list Servers from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		listServerExample, true)
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerGet, "get", "Get a Server",
		"Use this command to get information about a specified Server from a Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		getServerExample, true)
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), get.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, serverCmd, PreRunDataCenterId, RunServerCreate, "create", "Create a Server",
		`Use this command to create a Server in a specified Virtual Data Center. The name, cores, ram, cpu-family and availability zone options can be set.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, createServerExample, true)
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
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
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Server to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerUpdate, "update", "Update a Server",
		`Use this command to update a specified Server from a Virtual Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, updateServerExample, true)
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), update.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Server to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerDelete, "delete", "Delete a Server",
		`Use this command to delete a specified Server from a Virtual Data Center.

NOTE: This will not automatically remove the storage Volume(s) attached to a Server.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, deleteServerExample, true)
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), deleteCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Server to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be deleted [seconds]")

	/*
		Start Command
	*/
	start := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerStart, "start", "Start a Server",
		`Use this command to start a Server from a Virtual Data Center. If the Server's public IP was deallocated then a new IP will be assigned.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, startServerExample, true)
	start.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = start.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), start.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	start.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Server to start")
	start.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be started [seconds]")

	/*
		Stop Command
	*/
	stop := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerStop, "stop", "Stop a Server",
		`Use this command to stop a Server from a Virtual Data Center. The machine will be forcefully powered off, billing will cease, and the public IP, if one is allocated, will be deallocated.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, stopServerExample, true)
	stop.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = stop.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), stop.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	stop.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Server to stop")
	stop.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be stopped [seconds]")

	/*
		Reboot Command
	*/
	reboot := builder.NewCommand(ctx, serverCmd, PreRunDcServerIds, RunServerReboot, "reboot", "Force a hard reboot of a Server",
		`Use this command to force a hard reboot of the Server. Do not use this method if you want to gracefully reboot the machine. This is the equivalent of powering off the machine and turning it back on.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id`, resetServerExample, true)
	reboot.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	reboot.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = reboot.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverCmd.Name(), reboot.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	serverCmd.AddCommand(serverVolume())

	return serverCmd
}

func PreRunDcServerIds(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgServerId)
}

func RunServerList(c *builder.CommandConfig) error {
	servers, _, err := c.Servers().List(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)))
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
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
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
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerZone)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerCores)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerRAM)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:     server,
		KeyValue:       getServersKVMaps([]resources.Server{*server}),
		Columns:        getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "create",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
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
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		KeyValue:       getServersKVMaps([]resources.Server{*server}),
		Columns:        getServersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		OutputJSON:     server,
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "update",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunServerDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete server"); err != nil {
		return err
	}
	resp, err := c.Servers().Delete(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "delete",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunServerStart(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "start server"); err != nil {
		return err
	}
	resp, err := c.Servers().Start(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "start",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunServerStop(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "stop server"); err != nil {
		return err
	}
	resp, err := c.Servers().Stop(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "stop",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
}

func RunServerReboot(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "reboot server"); err != nil {
		return err
	}
	resp, err := c.Servers().Reboot(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse:    resp,
		Resource:       "server",
		Verb:           "reboot",
		WaitForRequest: viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest)),
	})
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
