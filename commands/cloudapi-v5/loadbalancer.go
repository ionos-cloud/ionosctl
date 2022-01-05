package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadBalancerCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultLoadbalancerCols, printer.ColsMessage(allLoadbalancerCols))
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
		LongDesc:   "Use this command to retrieve a list of Load Balancers within a Virtual Data Center on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LoadbalancersFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listLoadbalancerExample,
		PreCmdRun:  PreRunLoadBalancerList,
		CmdRun:     RunLoadBalancerList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadBalancersFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadBalancersFilters(), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgLoadBalancerId, cloudapiv5.ArgIdShort, "", cloudapiv5.LoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	create.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "Load Balancer", "Name of the Load Balancer")
	create.AddBoolFlag(cloudapiv5.ArgDhcp, "", cloudapiv5.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
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
	update.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgLoadBalancerId, cloudapiv5.ArgIdShort, "", cloudapiv5.LoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "Name of the Load Balancer")
	update.AddStringFlag(cloudapiv5.ArgIp, "", "", "The IP of the Load Balancer")
	update.AddBoolFlag(cloudapiv5.ArgDhcp, "", cloudapiv5.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
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
		PreCmdRun:  PreRunDcLoadBalancerDelete,
		CmdRun:     RunLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgLoadBalancerId, cloudapiv5.ArgIdShort, "", cloudapiv5.LoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Load Balancer deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Loadblancers from a virtual Datacenter.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer deletion [seconds]")

	loadbalancerCmd.AddCommand(LoadBalancerNicCmd())

	return loadbalancerCmd
}

func PreRunLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.LoadBalancersFilters(), completer.LoadbalancersFiltersUsage())
	}
	return nil
}

func PreRunDcLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgLoadBalancerId)
}

func PreRunDcLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgLoadBalancerId},
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgAll},
	)
}

func RunLoadBalancerList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting LoadBalancers from Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	lbs, resp, err := c.CloudApiV5Services.Loadbalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(nil, c, getLoadbalancers(lbs)))
}

func RunLoadBalancerGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting LoadBalancer with ID: %v from Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	lb, resp, err := c.CloudApiV5Services.Loadbalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLoadBalancerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(nil, c, []resources.Loadbalancer{*lb}))
}

func RunLoadBalancerCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgDhcp))
	c.Printer.Verbose("Properties set for creating the load balancer: Name: %v, Dhcp: %v", name, dhcp)
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	lb, resp, err := c.CloudApiV5Services.Loadbalancers().Create(dcId, name, dhcp)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getLoadbalancerPrint(resp, c, []resources.Loadbalancer{*lb}))
}

func RunLoadBalancerUpdate(c *core.CommandConfig) error {
	input := resources.LoadbalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgIp)) {
		ip := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIp))
		input.SetIp(ip)
		c.Printer.Verbose("Property Ip set: %v", ip)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgDhcp))
		input.SetDhcp(dhcp)
		c.Printer.Verbose("Property Dhcp set: %v", dhcp)
	}
	c.Printer.Verbose("Updating LoadBalancer with ID: %v from Datacenter with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	lb, resp, err := c.CloudApiV5Services.Loadbalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLoadBalancerId)),
		input,
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
	return c.Printer.Print(getLoadbalancerPrint(resp, c, []resources.Loadbalancer{*lb}))
}

func RunLoadBalancerDelete(c *core.CommandConfig) error {
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	loadBlanacerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLoadBalancerId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll)) {
		if err := DeleteAllLoadBalancers(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete loadbalancer"); err != nil {
			return err
		}
		c.Printer.Verbose("Datacenter ID: %v", dcid)
		c.Printer.Verbose("Starting deleting Load balancer with id: %v...", loadBlanacerId)
		resp, err := c.CloudApiV5Services.Loadbalancers().Delete(dcid, loadBlanacerId)
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getLoadbalancerPrint(resp, c, nil))
	}
}

func DeleteAllLoadBalancers(c *core.CommandConfig) error {
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	_ = c.Printer.Print("LoadBalancers to be deleted:")
	loadBalancers, _, err := c.CloudApiV5Services.Loadbalancers().List(dcid, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if loadBalancersItems, ok := loadBalancers.GetItemsOk(); ok && loadBalancersItems != nil {
		for _, lb := range *loadBalancersItems {
			toPrint := ""
			if id, ok := lb.GetIdOk(); ok && id != nil {
				toPrint += "LoadBalancer Id: " + *id
			}
			if properties, ok := lb.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					toPrint += " LoadBalancer Name: " + *name
				}
			}
			_ = c.Printer.Print(toPrint)
		}
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the LoadBalancers"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the LoadBalancers...")
		var multiErr error
		for _, lb := range *loadBalancersItems {
			if id, ok := lb.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Datacenter ID: %v", dcid)
				c.Printer.Verbose("Starting deleting Load balancer with id: %v...", *id)
				resp, err := c.CloudApiV5Services.Loadbalancers().Delete(dcid, *id)
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
		return errors.New("could not get items of LoadBalancers")
	}
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

func getLoadbalancerPrint(resp *resources.Response, c *core.CommandConfig, lbs []resources.Loadbalancer) printer.Result {
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

func getLoadbalancers(loadbalancers resources.Loadbalancers) []resources.Loadbalancer {
	vs := make([]resources.Loadbalancer, 0)
	if items, ok := loadbalancers.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			vs = append(vs, resources.Loadbalancer{Loadbalancer: s})
		}
	}
	return vs
}

func getLoadbalancersKVMaps(vs []resources.Loadbalancer) []map[string]interface{} {
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
