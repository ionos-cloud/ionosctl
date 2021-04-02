package commands

import (
	"context"
	"errors"
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

func loadBalancer() *builder.Command {
	loadbalancerCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "loadbalancer",
			Aliases:          []string{"lb"},
			Short:            "Load Balancer Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl loadbalancer` + "`" + ` manage your Load Balancers on your account. With the ` + "`" + `ionosctl loadbalancer` + "`" + ` command, you can list, create, delete Load Balancers and manage their configuration details.`,
			TraverseChildren: true,
		},
	}
	globalFlags := loadbalancerCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	viper.BindPFlag(builder.GetGlobalFlagName(loadbalancerCmd.Command.Use, config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	loadbalancerCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	viper.BindPFlag(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdValidate, RunLoadBalancerList, "list", "List Load Balancers",
		"Use this command to list all Load Balancers from a Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		listLoadbalancerExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdLoadBalancerIdValidate, RunLoadBalancerGet, "get", "Get a Load Balancer",
		"Use this command to retrieve information about a Load Balancer instance.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		getLoadbalancerExample, true)
	get.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	get.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, loadbalancerCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdValidate, RunLoadBalancerCreate, "create", "Create a Load Balancer",
		`Use this command to create a new Load Balancer on your account. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`, createLoadbalancerExample, true)
	create.AddStringFlag(config.ArgLoadBalancerName, "", "", "Name of the Load Balancer")
	create.AddBoolFlag(config.ArgLoadBalancerDhcp, "", config.DefaultLoadBalancerDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP")
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Load Balancer to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Load Balancer to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdLoadBalancerIdValidate, RunLoadBalancerUpdate, "update", "Update a Load Balancer",
		`Use this command to update the configuration of a specified Load Balancer.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`, updateLoadbalancerExample, true)
	update.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	update.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, loadbalancerCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgLoadBalancerName, "", "", "Name of the Load Balancer")
	update.AddStringFlag(config.ArgLoadBalancerIp, "", "", "The IP of the Load Balancer")
	update.AddBoolFlag(config.ArgLoadBalancerDhcp, "", config.DefaultLoadBalancerDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Load Balancer to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Load Balancer to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdLoadBalancerIdValidate, RunLoadBalancerDelete, "delete", "Delete a Load Balancer",
		`Use this command to permanently delete the specified Load Balancer. This action is irreversible.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`, deleteLoadbalancerExample, true)
	deleteCmd.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, loadbalancerCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Load Balancer to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Load Balancer to be deleted [seconds]")

	return loadbalancerCmd
}

func PreRunGlobalDcIdLoadBalancerIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLoadBalancerId)
	if err != nil {
		return err
	}
	return nil
}

func RunLoadBalancerList(c *builder.CommandConfig) error {
	loadbalancers, _, err := c.Loadbalancers().List(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	ss := getLoadbalancers(loadbalancers)
	return c.Printer.Print(printer.Result{
		OutputJSON: loadbalancers,
		KeyValue:   getLoadbalancersKVMaps(ss),
		Columns:    getLoadbalancersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunLoadBalancerGet(c *builder.CommandConfig) error {
	loadbalancer, _, err := c.Loadbalancers().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: loadbalancer,
		KeyValue:   getLoadbalancersKVMaps([]resources.Loadbalancer{*loadbalancer}),
		Columns:    getLoadbalancersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunLoadBalancerCreate(c *builder.CommandConfig) error {
	loadbalancer, resp, err := c.Loadbalancers().Create(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerName)),
		viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerDhcp)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  loadbalancer,
		KeyValue:    getLoadbalancersKVMaps([]resources.Loadbalancer{*loadbalancer}),
		Columns:     getLoadbalancersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "loadbalancer",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunLoadBalancerUpdate(c *builder.CommandConfig) error {
	input := resources.LoadbalancerProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerIp)) {
		input.SetIp(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerIp)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerDhcp)) {
		input.SetDhcp(viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerDhcp)))
	}
	loadbalancer, resp, err := c.Loadbalancers().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
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
		OutputJSON:  loadbalancer,
		KeyValue:    getLoadbalancersKVMaps([]resources.Loadbalancer{*loadbalancer}),
		Columns:     getLoadbalancersCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "loadbalancer",
		Verb:        "update",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunLoadBalancerDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete loadbalancer")
	if err != nil {
		return err
	}
	resp, err := c.Loadbalancers().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
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
		Resource:    "loadbalancer",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultLoadbalancerCols = []string{"LoadBalancerId", "Name", "Dhcp"}

type LoadbalancerPrint struct {
	LoadbalancerId string `json:"LoadBalancerId,omitempty"`
	Name           string `json:"Name,omitempty"`
	Dhcp           bool   `json:"Dhcp,omitempty"`
	Ip             string `json:"Ip,omitempty"`
}

func getLoadbalancersCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultLoadbalancerCols
	}

	columnsMap := map[string]string{
		"LoadBalancerId": "LoadBalancerId",
		"Name":           "Name",
		"Dhcp":           "Dhcp",
		"Ip":             "Ip",
	}
	var loadbalancerCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			loadbalancerCols = append(loadbalancerCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return loadbalancerCols
}

func getLoadbalancers(loadbalancers resources.Loadbalancers) []resources.Loadbalancer {
	vs := make([]resources.Loadbalancer, 0)
	for _, s := range *loadbalancers.Items {
		vs = append(vs, resources.Loadbalancer{s})
	}
	return vs
}

func getLoadbalancersKVMaps(vs []resources.Loadbalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(vs))
	for _, v := range vs {
		properties := v.GetProperties()
		var loadbalancerPrint LoadbalancerPrint
		if id, ok := v.GetIdOk(); ok && id != nil {
			loadbalancerPrint.LoadbalancerId = *id
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			loadbalancerPrint.Name = *name
		}
		if dhcp, ok := properties.GetDhcpOk(); ok && dhcp != nil {
			loadbalancerPrint.Dhcp = *dhcp
		}
		if ip, ok := properties.GetIpOk(); ok && ip != nil {
			loadbalancerPrint.Ip = *ip
		}
		o := structs.Map(loadbalancerPrint)
		out = append(out, o)
	}
	return out
}

func getLoadbalancersIds(outErr io.Writer, parentCmdName string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	loadbalancerSvc := resources.NewLoadbalancerService(clientSvc.Get(), context.TODO())
	loadbalancers, _, err := loadbalancerSvc.List(viper.GetString(builder.GetGlobalFlagName(parentCmdName, config.ArgDataCenterId)))
	clierror.CheckError(err, outErr)
	loadbalancersIds := make([]string, 0)
	if loadbalancers.Loadbalancers.Items != nil {
		for _, v := range *loadbalancers.Loadbalancers.Items {
			loadbalancersIds = append(loadbalancersIds, *v.GetId())
		}
	} else {
		return nil
	}
	return loadbalancersIds
}
