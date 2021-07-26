package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loadBalancer() *core.Command {
	ctx := context.TODO()
	loadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "loadbalancer",
			Aliases:          []string{"lb"},
			Short:            "Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl loadbalancer` manage your Load Balancers on your account. With the `ionosctl loadbalancer` command, you can list, create, delete Load Balancers and manage their configuration details.",
			TraverseChildren: true,
		},
	}
	globalFlags := loadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultLoadbalancerCols, utils.ColsMessage(allLoadbalancerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(loadbalancerCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = loadbalancerCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLoadbalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "loadbalancer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Load Balancers",
		LongDesc:   "Use this command to retrieve a list of Load Balancers within a Virtual Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLoadbalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLoadBalancerList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "loadbalancer",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Load Balancer",
		LongDesc:   "Use this command to retrieve information about a Load Balancer instance.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		Example:    getLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagLoadBalancerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Load Balancer",
		LongDesc: `Use this command to create a new Load Balancer within the Virtual Data Center. Load balancers can be used for public or private IP traffic. The name, IP and DHCP for the Load Balancer can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createLoadbalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLoadBalancerCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Load Balancer", "Name of the Load Balancer")
	create.AddBoolFlag(config.ArgDhcp, "", config.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Load Balancer creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Load Balancer",
		LongDesc: `Use this command to update the configuration of a specified Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`,
		Example:    updateLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagLoadBalancerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Load Balancer")
	update.AddStringFlag(config.ArgIp, "", "", "The IP of the Load Balancer")
	update.AddBoolFlag(config.ArgDhcp, "", config.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Load Balancer update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, loadbalancerCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Load Balancer",
		LongDesc: `Use this command to delete the specified Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`,
		Example:    deleteLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagLoadBalancerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Load Balancer deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer deletion [seconds]")

	loadbalancerCmd.AddCommand(loadBalancerNic())

	return loadbalancerCmd
}

func PreRunDcLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgLoadBalancerId)
}

func RunLoadBalancerList(c *core.CommandConfig) error {
	lbs, _, err := c.Loadbalancers().List(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(nil, c, getLoadbalancers(lbs)))
}

func RunLoadBalancerGet(c *core.CommandConfig) error {
	lb, _, err := c.Loadbalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(nil, c, []v5.Loadbalancer{*lb}))
}

func RunLoadBalancerCreate(c *core.CommandConfig) error {
	lb, resp, err := c.Loadbalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgName)),
		viper.GetBool(core.GetFlagName(c.NS, config.ArgDhcp)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(resp, c, []v5.Loadbalancer{*lb}))
}

func RunLoadBalancerUpdate(c *core.CommandConfig) error {
	input := v5.LoadbalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIp)) {
		input.SetIp(viper.GetString(core.GetFlagName(c.NS, config.ArgIp)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDhcp)) {
		input.SetDhcp(viper.GetBool(core.GetFlagName(c.NS, config.ArgDhcp)))
	}
	lb, resp, err := c.Loadbalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(resp, c, []v5.Loadbalancer{*lb}))
}

func RunLoadBalancerDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete loadbalancer"); err != nil {
		return err
	}
	resp, err := c.Loadbalancers().Delete(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(resp, c, nil))
}

// Output Printing

var (
	defaultLoadbalancerCols = []string{"LoadBalancerId", "Name", "Dhcp", "State"}
	allLoadbalancerCols     = []string{"LoadBalancerId", "Name", "Dhcp", "State", "Ip"}
)

type LoadbalancerPrint struct {
	LoadBalancerId string `json:"LoadBalancerId,omitempty"`
	Name           string `json:"Name,omitempty"`
	Dhcp           bool   `json:"Dhcp,omitempty"`
	Ip             string `json:"Ip,omitempty"`
	State          string `json:"State,omitempty"`
}

func getLoadbalancerPrint(resp *v5.Response, c *core.CommandConfig, lbs []v5.Loadbalancer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if lbs != nil {
			r.OutputJSON = lbs
			r.KeyValue = getLoadbalancersKVMaps(lbs)
			r.Columns = getLoadbalancersCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
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
		"State":          "State",
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

func getLoadbalancers(loadbalancers v5.Loadbalancers) []v5.Loadbalancer {
	vs := make([]v5.Loadbalancer, 0)
	for _, s := range *loadbalancers.Items {
		vs = append(vs, v5.Loadbalancer{Loadbalancer: s})
	}
	return vs
}

func getLoadbalancersKVMaps(vs []v5.Loadbalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(vs))
	for _, v := range vs {
		var loadbalancerPrint LoadbalancerPrint
		if id, ok := v.GetIdOk(); ok && id != nil {
			loadbalancerPrint.LoadBalancerId = *id
		}
		if properties, ok := v.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				loadbalancerPrint.Name = *name
			}
			if dhcp, ok := properties.GetDhcpOk(); ok && dhcp != nil {
				loadbalancerPrint.Dhcp = *dhcp
			}
			if ip, ok := properties.GetIpOk(); ok && ip != nil {
				loadbalancerPrint.Ip = *ip
			}
		}
		if metadata, ok := v.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				loadbalancerPrint.State = *state
			}
		}
		o := structs.Map(loadbalancerPrint)
		out = append(out, o)
	}
	return out
}

func getLoadbalancersIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	loadbalancerSvc := v5.NewLoadbalancerService(clientSvc.Get(), context.TODO())
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
