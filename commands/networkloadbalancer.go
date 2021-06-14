package commands

import (
	"context"
	"errors"
	"io"
	"os"

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

func networkloadbalancer() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultNetworkLoadBalancerCols, utils.ColsMessage(defaultNetworkLoadBalancerCols))
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
		LongDesc:   "Use this command to list Network Load Balancers from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listNetworkLoadBalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunNetworkLoadBalancerList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgNetworkLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagNetworkLoadBalancerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Network Load Balancer to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for waiting for Network Load Balancer to be in AVAILABLE state [seconds]")

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
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	create.AddIntFlag(config.ArgListenerLan, "", 2, "Id of the listening LAN")
	create.AddIntFlag(config.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	create.AddStringSliceFlag(config.ArgIps, "", []string{""}, "Collection of IP addresses of the Network Load Balancer")
	create.AddStringSliceFlag(config.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer creation [seconds]")

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
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgNetworkLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagNetworkLoadBalancerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	update.AddIntFlag(config.ArgListenerLan, "", 2, "Id of the listening LAN")
	update.AddIntFlag(config.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	update.AddStringSliceFlag(config.ArgIps, "", []string{""}, "Collection of IP addresses of the Network Load Balancer")
	update.AddStringSliceFlag(config.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer update [seconds]")

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
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgNetworkLoadBalancerId, config.ArgIdShort, "", config.RequiredFlagNetworkLoadBalancerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Network Load Balancer deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Network Load Balancer deletion [seconds]")

	networkloadbalancerCmd.AddCommand(networkloadbalancerFlowLog())
	networkloadbalancerCmd.AddCommand(networkloadbalancerRule())

	return networkloadbalancerCmd
}

func PreRunDcNetworkLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNetworkLoadBalancerId)
}

func RunNetworkLoadBalancerList(c *core.CommandConfig) error {
	networkloadbalancers, _, err := c.NetworkLoadBalancers().List(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(nil, c, getNetworkLoadBalancers(networkloadbalancers)))
}

func RunNetworkLoadBalancerGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, GetStateNetworkLoadBalancer, viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId))); err != nil {
		return err
	}
	ng, _, err := c.NetworkLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(nil, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerCreate(c *core.CommandConfig) error {
	proper := getNewNetworkLoadBalancerInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, config.ArgTargetLan)))
	}
	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, config.ArgListenerLan)))
	}
	ng, resp, err := c.NetworkLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		resources.NetworkLoadBalancer{
			NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
				Properties: &proper.NetworkLoadBalancerProperties,
			},
		},
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerUpdate(c *core.CommandConfig) error {
	input := getNewNetworkLoadBalancerInfo(c)
	ng, resp, err := c.NetworkLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		*input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, []resources.NetworkLoadBalancer{*ng}))
}

func RunNetworkLoadBalancerDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete network load balancer"); err != nil {
		return err
	}
	resp, err := c.NetworkLoadBalancers().Delete(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNetworkLoadBalancerPrint(resp, c, nil))
}

func getNewNetworkLoadBalancerInfo(c *core.CommandConfig) *resources.NetworkLoadBalancerProperties {
	input := ionoscloud.NetworkLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIps)) {
		input.SetIps(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgIps)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgListenerLan)) {
		input.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, config.ArgListenerLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgTargetLan)) {
		input.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, config.ArgTargetLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPrivateIps)) {
		input.SetLbPrivateIps(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgPrivateIps)))
	}
	return &resources.NetworkLoadBalancerProperties{
		NetworkLoadBalancerProperties: input,
	}
}

// Wait for State

func GetStateNetworkLoadBalancer(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.NetworkLoadBalancers().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)), objId)
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
	ss := make([]resources.NetworkLoadBalancer, 0)
	for _, s := range *networkloadbalancers.Items {
		ss = append(ss, resources.NetworkLoadBalancer{NetworkLoadBalancer: s})
	}
	return ss
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

func getNetworkLoadBalancersIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(clientSvc.Get(), context.TODO())
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := networkloadbalancers.NetworkLoadBalancers.GetItemsOk(); ok && items != nil {
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
