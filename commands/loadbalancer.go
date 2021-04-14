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
	_ = viper.BindPFlag(builder.GetGlobalFlagName(loadbalancerCmd.Command.Use, config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = loadbalancerCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

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
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Load Balancer to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Load Balancer to be deleted [seconds]")

	/*
		Attach Nic Command
	*/
	attachNic := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdNicLoadBalancerIdsValidate, RunLoadBalancerAttachNic, "attach-nic", "Attach a NIC to a Load Balancer",
		`Use this command to attach a specified NIC to a Load Balancer on your account.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`, attachNicLoadbalancerExample, true)
	attachNic.AddStringFlag(config.ArgServerId, "", "", "The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion")
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(loadbalancerCmd.Command.Name(), attachNic.Command.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to attach to a Load Balancer")
	attachNic.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be attached to a Load Balancer [seconds]")

	/*
		List Nics Command
	*/
	listAttached := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdLoadBalancerIdValidate, RunLoadBalancerListNics, "list-nics", "List attached NICs from a Load Balancer",
		"Use this command to get a list of attached NICs to a Load Balancer from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		listNicsLoadbalancerExample, true)
	listAttached.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = listAttached.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Nic Command
	*/
	getAttached := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdNicLoadBalancerIdsValidate, RunLoadBalancerGetNic, "get-nic", "Get an attached NIC to a Load Balancer",
		"Use this command to retrieve information about an attached NIC to a Load Balancer.\n\nRequired values to run the command:\n\n* Data Center Id\n* Load Balancer Id\n* NIC Id",
		getNicLoadbalancerExample, true)
	getAttached.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getAttached.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(
			os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetFlagName(loadbalancerCmd.Command.Name(), getAttached.Command.Name(), config.ArgLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Nic Command
	*/
	detachNic := builder.NewCommand(context.TODO(), loadbalancerCmd, PreRunGlobalDcIdNicLoadBalancerIdsValidate, RunLoadBalancerDetachNic, "detach-nic", "Detach a NIC from a Load Balancer",
		`Use this command to detach a NIC from a Load Balancer.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`, detachNicLoadbalancerExample, true)
	detachNic.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(loadbalancerCmd.Command.Name(), detachNic.Command.Name(), config.ArgLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(loadbalancerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to detach from a Load Balancer")
	detachNic.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be detached from a Load Balancer [seconds]")

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

func PreRunGlobalDcIdNicLoadBalancerIdsValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgNicId, config.ArgLoadBalancerId)
}

func RunLoadBalancerAttachNic(c *builder.CommandConfig) error {
	attachedNic, resp, err := c.Loadbalancers().AttachNic(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, getNic(attachedNic)))
}

func RunLoadBalancerListNics(c *builder.CommandConfig) error {
	attachedNics, _, err := c.Loadbalancers().ListNics(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getAttachedNics(attachedNics)))
}

func RunLoadBalancerGetNic(c *builder.CommandConfig) error {
	n, _, err := c.Loadbalancers().GetNic(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getNic(n)))
}

func RunLoadBalancerDetachNic(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "detach nic from loadbalancer")
	if err != nil {
		return err
	}
	resp, err := c.Loadbalancers().DetachNic(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, nil))
}

// Output Printing

var defaultLoadbalancerCols = []string{"LoadBalancerId", "Name", "Dhcp"}

type LoadbalancerPrint struct {
	LoadBalancerId string `json:"LoadBalancerId,omitempty"`
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
		vs = append(vs, resources.Loadbalancer{Loadbalancer: s})
	}
	return vs
}

func getLoadbalancersKVMaps(vs []resources.Loadbalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(vs))
	for _, v := range vs {
		properties := v.GetProperties()
		var loadbalancerPrint LoadbalancerPrint
		if id, ok := v.GetIdOk(); ok && id != nil {
			loadbalancerPrint.LoadBalancerId = *id
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

func getLoadbalancersIds(outErr io.Writer, datacenterId string) []string {
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
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId)
	clierror.CheckError(err, outErr)
	loadbalancersIds := make([]string, 0)
	if items, ok := loadbalancers.Loadbalancers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				loadbalancersIds = append(loadbalancersIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return loadbalancersIds
}

func getAttachedNicsIds(outErr io.Writer, datacenterId, loadbalancerId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	nicSvc := resources.NewLoadbalancerService(clientSvc.Get(), context.TODO())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId)
	clierror.CheckError(err, outErr)
	attachedNicsIds := make([]string, 0)
	if items, ok := nics.BalancedNics.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedNicsIds = append(attachedNicsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedNicsIds
}
