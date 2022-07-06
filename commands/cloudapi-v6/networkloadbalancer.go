package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkloadbalancerCmd() *core.Command {
	ctx := context.TODO()
	networkloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "networkloadbalancer",
			Aliases:          []string{"nlb"},
			Short:            "Network Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer` allow you to create, list, get, update, delete Network Load Balancers.",
			TraverseChildren: true,
		},
	}
	globalFlags := networkloadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultNetworkLoadBalancerCols, printer.ColsMessage(defaultNetworkLoadBalancerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(networkloadbalancerCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = networkloadbalancerCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNetworkLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, networkloadbalancerCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "networkloadbalancer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancers",
		LongDesc:   "Use this command to list Network Load Balancers from a specified Virtual Data Center.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.NlbsFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listNetworkLoadBalancerExample,
		PreCmdRun:  PreRunNetworkLoadBalancerList,
		CmdRun:     RunNetworkLoadBalancerList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultListDepth, "Controls the detail depth of the response objects. Max depth is 10.")
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, networkloadbalancerCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "networkloadbalancer",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Network Load Balancer",
		LongDesc:   "Use this command to get information about a specified Network Load Balancer from a Virtual Data Center. You can also wait for Network Load Balancer to get in AVAILABLE state using `--wait-for-state` option.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id",
		Example:    getNetworkLoadBalancerExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Network Load Balancer to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for waiting for Network Load Balancer to be in AVAILABLE state [seconds]")
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	get.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultGetDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, networkloadbalancerCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer",
		LongDesc: `Use this command to create a Network Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createNetworkLoadBalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunNetworkLoadBalancerCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	create.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	create.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	create.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of IP addresses of the Network Load Balancer")
	create.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer creation [seconds]")
	create.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultCreateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, networkloadbalancerCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer",
		LongDesc: `Use this command to update a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id`,
		Example:    updateNetworkLoadBalancerExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	update.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	update.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of IP addresses of the Network Load Balancer")
	update.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer update [seconds]")
	update.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultUpdateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, networkloadbalancerCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Network Load Balancer",
		LongDesc: `Use this command to delete a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id`,
		Example:    deleteNetworkLoadBalancerExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerDelete,
		CmdRun:     RunNetworkLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all NetworkLoadBalancers.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer deletion [seconds]")
	deleteCmd.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultDeleteDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	networkloadbalancerCmd.AddCommand(NetworkloadbalancerFlowLogCmd())
	networkloadbalancerCmd.AddCommand(NetworkloadbalancerRuleCmd())

	return networkloadbalancerCmd
}

func PreRunNetworkLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.NlbsFilters(), completer.NlbsFiltersUsage())
	}
	return nil
}

func PreRunDcNetworkLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId)
}

func PreRunDcNetworkLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunNetworkLoadBalancerList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("List Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
		if !structs.IsZero(listQueryParams.QueryParams) {
			c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams.QueryParams))
		}
	}
	networkloadbalancers, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(nil, c, getNetworkLoadBalancers(networkloadbalancers)))
}

func RunNetworkLoadBalancerGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	if err := utils.WaitForState(c, waiter.NetworkLoadBalancerStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))); err != nil {
		return err
	}
	c.Printer.Verbose("NetworkLoadBalancer with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)))
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(nil, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	proper := getNewNetworkLoadBalancerInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		resources.NetworkLoadBalancer{
			NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
				Properties: &proper.NetworkLoadBalancerProperties,
			},
		},
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	input := getNewNetworkLoadBalancerInfo(c)
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		*input,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNetworkLoadBalancers(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete network load balancer"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting NetworkLoadBalancer with id: %v...", nlbId)
		resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Delete(dcId, nlbId)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, nil))
	}
}

func getNewNetworkLoadBalancerInfo(c *core.CommandConfig) *resources.NetworkLoadBalancerProperties {
	input := ionoscloud.NetworkLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		ips := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
		input.SetIps(ips)
		c.Printer.Verbose("Property Ips set: %v", ips)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)) {
		listenerLan := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan))
		input.SetListenerLan(listenerLan)
		c.Printer.Verbose("Property ListenerLan set: %v", listenerLan)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)) {
		targetLan := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan))
		input.SetTargetLan(targetLan)
		c.Printer.Verbose("Property TargetLan set: %v", targetLan)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)) {
		privateIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps))
		input.SetLbPrivateIps(privateIps)
		c.Printer.Verbose("Property PrivateIps set: %v", privateIps)

	}
	return &resources.NetworkLoadBalancerProperties{
		NetworkLoadBalancerProperties: input,
	}
}

func DeleteAllNetworkLoadBalancers(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting NetworkLoadBalancers...")
	networkLoadBalancers, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().List(dcId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if nlbItems, ok := networkLoadBalancers.GetItemsOk(); ok && nlbItems != nil {
		if len(*nlbItems) > 0 {
			_ = c.Printer.Print("NetworkLoadBalancers to be deleted:")
			for _, networkLoadBalancer := range *nlbItems {
				toPrint := ""
				if id, ok := networkLoadBalancer.GetIdOk(); ok && id != nil {
					toPrint += "NetworkLoadBalancer Id: " + *id
				}
				if properties, ok := networkLoadBalancer.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += " NetworkLoadBalancer Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the NetworkLoadBalancers"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the NetworkLoadBalancers...")
			var multiErr error
			for _, networkLoadBalancer := range *nlbItems {
				if id, ok := networkLoadBalancer.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting NetworkLoadBalancer with id: %v...", *id)
					resp, err = c.CloudApiV6Services.NetworkLoadBalancers().Delete(dcId, *id)
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
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no NetworkLoadBalancers found")
		}
	} else {
		return errors.New("could not get items of NetworkLoadBalancers")
	}
}

// Output Printing

var defaultNetworkLoadBalancerCols = []string{"NetworkLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "LbPrivateIps", "State"}

type NetworkLoadBalancerPrint struct {
	NetworkLoadBalancerId string   `json:"NetworkLoadBalancerId,omitempty"`
	Name                  string   `json:"Name,omitempty"`
	ListenerLan           int32    `json:"ListenerLan,omitempty"`
	Ips                   []string `json:"Ips,omitempty"`
	TargetLan             int32    `json:"TargetLan,omitempty"`
	LbPrivateIps          []string `json:"LbPrivateIps,omitempty"`
	State                 string   `json:"State,omitempty"`
}

func getNetworkLoadBalancerPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NetworkLoadBalancer) printer.Result {
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
			r.KeyValue = getNetworkLoadBalancersKVMaps(ss)
			r.Columns = getNetworkLoadBalancersCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getNetworkLoadBalancersCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultNetworkLoadBalancerCols
	}

	columnsMap := map[string]string{
		"NetworkLoadBalancerId": "NetworkLoadBalancerId",
		"Name":                  "Name",
		"ListenerLan":           "ListenerLan",
		"Ips":                   "Ips",
		"TargetLan":             "TargetLan",
		"LbPrivateIps":          "LbPrivateIps",
		"State":                 "State",
	}
	var networkloadbalancerCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			networkloadbalancerCols = append(networkloadbalancerCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return networkloadbalancerCols
}

func getNetworkLoadBalancers(networkloadbalancers resources.NetworkLoadBalancers) []resources.NetworkLoadBalancer {
	networkLoadBalancerObjs := make([]resources.NetworkLoadBalancer, 0)
	if items, ok := networkloadbalancers.GetItemsOk(); ok && items != nil {
		for _, networkLoadBalancer := range *items {
			networkLoadBalancerObjs = append(networkLoadBalancerObjs, resources.NetworkLoadBalancer{NetworkLoadBalancer: networkLoadBalancer})
		}
	}
	return networkLoadBalancerObjs
}

func getNetworkLoadBalancersKVMaps(ss []resources.NetworkLoadBalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var networkloadbalancerPrint NetworkLoadBalancerPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			networkloadbalancerPrint.NetworkLoadBalancerId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				networkloadbalancerPrint.Name = *name
			}
			if listenerLan, ok := properties.GetListenerLanOk(); ok && listenerLan != nil {
				networkloadbalancerPrint.ListenerLan = *listenerLan
			}
			if ips, ok := properties.GetIpsOk(); ok && ips != nil {
				networkloadbalancerPrint.Ips = *ips
			}
			if targetLan, ok := properties.GetTargetLanOk(); ok && targetLan != nil {
				networkloadbalancerPrint.TargetLan = *targetLan
			}
			if lbPrivateIps, ok := properties.GetLbPrivateIpsOk(); ok && lbPrivateIps != nil {
				networkloadbalancerPrint.LbPrivateIps = *lbPrivateIps
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				networkloadbalancerPrint.State = *state
			}
		}
		o := structs.Map(networkloadbalancerPrint)
		out = append(out, o)
	}
	return out
}
