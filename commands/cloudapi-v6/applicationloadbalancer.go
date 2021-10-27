package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerCmd() *core.Command {
	ctx := context.TODO()
	applicationloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "applicationloadbalancer",
			Aliases:          []string{"alb"},
			Short:            "Application Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl applicationloadbalancer` allow you to create, list, get, update, delete Application Load Balancers.",
			TraverseChildren: true,
		},
	}
	globalFlags := applicationloadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultApplicationLoadBalancerCols, printer.ColsMessage(defaultApplicationLoadBalancerCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(applicationloadbalancerCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = applicationloadbalancerCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultApplicationLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "applicationloadbalancer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Application Load Balancers",
		LongDesc:   "Use this command to list Application Load Balancers from a specified Virtual Data Center.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunApplicationLoadBalancerList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "applicationloadbalancer",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Application Load Balancer",
		LongDesc:   "Use this command to get information about a specified Application Load Balancer from a Virtual Data Center. You can also wait for Application Load Balancer to get in AVAILABLE state using `--wait-for-state` option.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id",
		Example:    getApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Application Load Balancer to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for waiting for Application Load Balancer to be in AVAILABLE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Application Load Balancer",
		LongDesc: `Use this command to create a Application Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunApplicationLoadBalancerCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Application Load Balancer", "Name of the Application Load Balancer")
	create.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	create.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	create.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of IP addresses of the Application Load Balancer. IP of the Listener Lan must be a customer reserved IP for the public Load Balancer and private IP for the private Load Balancer")
	create.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Application Load Balancer",
		LongDesc: `Use this command to update a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id`,
		Example:    updateApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Application Load Balancer", "Name of the Application Load Balancer")
	update.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	update.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	update.AddStringSliceFlag(cloudapiv6.ArgIps, "", []string{""}, "Collection of IP addresses of the Application Load Balancer. IP of the Listener Lan must be a customer reserved IP for the public Load Balancer and private IP for the private Load Balancer")
	update.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", []string{""}, "Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, applicationloadbalancerCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Application Load Balancer",
		LongDesc: `Use this command to delete a specified Application Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id`,
		Example:    deleteApplicationLoadBalancerExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Application Load Balancer deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Application Load Balancer deletion [seconds]")

	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerRuleCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerFlowLogCmd())

	return applicationloadbalancerCmd
}

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}

func RunApplicationLoadBalancerList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting ApplicationLoadBalancers from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationloadbalancers, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(nil, c, getApplicationLoadBalancers(applicationloadbalancers)))
}

func RunApplicationLoadBalancerGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if err := utils.WaitForState(c, GetStateApplicationLoadBalancer, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId))); err != nil {
		return err
	}
	ng, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(nil, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerCreate(c *core.CommandConfig) error {
	proper := getNewApplicationLoadBalancerInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Printer.Verbose("Property TargetLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Printer.Verbose("Property ListenerLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	c.Printer.Verbose("Creating ApplicationLoadBalancer in Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		resources.ApplicationLoadBalancer{
			ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
				Properties: &proper.ApplicationLoadBalancerProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerUpdate(c *core.CommandConfig) error {
	input := getNewApplicationLoadBalancerInfo(c)
	c.Printer.Verbose("Updating ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		*input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, []resources.ApplicationLoadBalancer{*ng}))
}

func RunApplicationLoadBalancerDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer"); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Delete(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getApplicationLoadBalancerPrint(resp, c, nil))
}

func getNewApplicationLoadBalancerInfo(c *core.CommandConfig) *resources.ApplicationLoadBalancerProperties {
	input := ionoscloud.ApplicationLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		input.SetIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
		c.Printer.Verbose("Property IPs set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)) {
		input.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Printer.Verbose("Property ListenerLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)) {
		input.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Printer.Verbose("Property TargetLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)) {
		input.SetLbPrivateIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
		c.Printer.Verbose("Property LbPrivateIps set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
	}
	return &resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: input,
	}
}

// Wait for State

func GetStateApplicationLoadBalancer(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), objId)
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

var defaultApplicationLoadBalancerCols = []string{"ApplicationLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "LbPrivateIps", "State"}

type ApplicationLoadBalancerPrint struct {
	ApplicationLoadBalancerId string   `json:"ApplicationLoadBalancerId,omitempty"`
	Name                      string   `json:"Name,omitempty"`
	ListenerLan               int32    `json:"ListenerLan,omitempty"`
	Ips                       []string `json:"Ips,omitempty"`
	TargetLan                 int32    `json:"TargetLan,omitempty"`
	LbPrivateIps              []string `json:"LbPrivateIps,omitempty"`
	State                     string   `json:"State,omitempty"`
}

func getApplicationLoadBalancerPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.ApplicationLoadBalancer) printer.Result {
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
			r.KeyValue = getApplicationLoadBalancersKVMaps(ss)
			r.Columns = getApplicationLoadBalancersCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getApplicationLoadBalancersCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultApplicationLoadBalancerCols
	}

	columnsMap := map[string]string{
		"ApplicationLoadBalancerId": "ApplicationLoadBalancerId",
		"Name":                      "Name",
		"ListenerLan":               "ListenerLan",
		"Ips":                       "Ips",
		"TargetLan":                 "TargetLan",
		"LbPrivateIps":              "LbPrivateIps",
		"State":                     "State",
	}
	var applicationloadbalancerCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			applicationloadbalancerCols = append(applicationloadbalancerCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return applicationloadbalancerCols
}

func getApplicationLoadBalancers(applicationloadbalancers resources.ApplicationLoadBalancers) []resources.ApplicationLoadBalancer {
	ss := make([]resources.ApplicationLoadBalancer, 0)
	for _, s := range *applicationloadbalancers.Items {
		ss = append(ss, resources.ApplicationLoadBalancer{ApplicationLoadBalancer: s})
	}
	return ss
}

func getApplicationLoadBalancersKVMaps(ss []resources.ApplicationLoadBalancer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var applicationloadbalancerPrint ApplicationLoadBalancerPrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			applicationloadbalancerPrint.ApplicationLoadBalancerId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				applicationloadbalancerPrint.Name = *name
			}
			if listenerLan, ok := properties.GetListenerLanOk(); ok && listenerLan != nil {
				applicationloadbalancerPrint.ListenerLan = *listenerLan
			}
			if ips, ok := properties.GetIpsOk(); ok && ips != nil {
				applicationloadbalancerPrint.Ips = *ips
			}
			if targetLan, ok := properties.GetTargetLanOk(); ok && targetLan != nil {
				applicationloadbalancerPrint.TargetLan = *targetLan
			}
			if lbPrivateIps, ok := properties.GetLbPrivateIpsOk(); ok && lbPrivateIps != nil {
				applicationloadbalancerPrint.LbPrivateIps = *lbPrivateIps
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				applicationloadbalancerPrint.State = *state
			}
		}
		o := structs.Map(applicationloadbalancerPrint)
		out = append(out, o)
	}
	return out
}
