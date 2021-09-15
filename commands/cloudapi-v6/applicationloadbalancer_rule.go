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

func ApplicationLoadBalancerRuleCmd() *core.Command {
	ctx := context.TODO()
	albRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Application Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl alb rule` allow you to create, list, get, update, delete Application Load Balancer Forwarding Rules.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, albRuleCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "rule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Application Load Balancer Forwarding Rules",
		LongDesc:   "Use this command to list Application Load Balancer Forwarding Rules from a specified Application Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id",
		Example:    listApplicationLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, albRuleCmd, core.CommandBuilder{
		Namespace:  "applicationloadbalancer",
		Resource:   "rule",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Application Load Balancer Forwarding Rule",
		LongDesc:   "Use this command to get information about a specified Application Load Balancer Forwarding Rule from a Application Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id\n* Forwarding Rule Id",
		Example:    getApplicationLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, albRuleCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a Application Load Balancer Forwarding Rule in a specified Application Load Balancer. You can also set Health Check Settings for Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Listener Ip
* Listener Port`,
		Example:    createApplicationLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunApplicationLoadBalancerForwardingRuleCreate,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name for the Forwarding Rule")
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Protocol of the balancing")
	create.AddStringFlag(cloudapiv6.ArgListenerIp, "", "", "Listening IP", core.RequiredFlagOption())
	create.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	create.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, albRuleCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to update a specified Application Load Balancer Forwarding Rule from a Application Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id`,
		Example:    updateApplicationLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Forwarding Rule")
	update.AddStringFlag(cloudapiv6.ArgListenerIp, "", "", "Listening IP")
	update.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening port number. Range: 1 to 65535")
	update.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, cloudapiv6.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, albRuleCmd, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to delete a specified Application Load Balancer Forwarding Rule from a Application Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id`,
		Example:    deleteApplicationLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule deletion [seconds]")

	albRuleCmd.AddCommand(AlbRuleHttpRuleCmd())

	return albRuleCmd
}

func PreRunApplicationLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgListenerIp, cloudapiv6.ArgListenerPort)
}

func PreRunDcApplicationLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId)
}

func RunApplicationLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting ForwardingRules for ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	albForwardingRules, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(nil, c, getAlbForwardingRules(albForwardingRules)))
}

func RunApplicationLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting ForwardingRule with ID: %v for ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationLoadBalancerForwardingRule, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(nil, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	proper := getAlbForwardingRulePropertiesSet(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	c.Printer.Verbose("Creating ForwardingRule for ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.ApplicationLoadBalancerForwardingRule{
			ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
				Properties: &proper.ApplicationLoadBalancerForwardingRuleProperties,
			},
		},
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	input := getAlbForwardingRulePropertiesSet(c)
	c.Printer.Verbose("Updating ForwardingRule with ID: %v from ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer forwarding rule"); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting ForwardingRule with ID: %v from ApplicationLoadBalancer with ID: %v in Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().DeleteForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, nil))
}

func getAlbForwardingRulePropertiesSet(c *core.CommandConfig) *resources.ApplicationLoadBalancerForwardingRuleProperties {
	input := ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)) {
		input.SetListenerIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)))
		c.Printer.Verbose("Property ListenerIp set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)) {
		input.SetListenerPort(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)))
		c.Printer.Verbose("Property ListenerPort set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)) {
		input.SetServerCertificates(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
		c.Printer.Verbose("Property ServerCertificates set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
	}
	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: input,
	}
}

// Output Printing

var defaultAlbForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Protocol", "ListenerIp", "ListenerPort", "ServerCertificates", "State"}

type AlbForwardingRulePrint struct {
	ForwardingRuleId   string   `json:"ForwardingRuleId,omitempty"`
	Name               string   `json:"Name,omitempty"`
	Protocol           string   `json:"Protocol,omitempty"`
	ListenerIp         string   `json:"ListenerIp,omitempty"`
	ListenerPort       int32    `json:"ListenerPort,omitempty"`
	ServerCertificates []string `json:"ServerCertificates,omitempty"`
	State              string   `json:"State,omitempty"`
}

func getAlbForwardingRulePrint(resp *resources.Response, c *core.CommandConfig, ss []resources.ApplicationLoadBalancerForwardingRule) printer.Result {
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
			r.KeyValue = getAlbForwardingRulesKVMaps(ss)
			r.Columns = getAlbForwardingRulesCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getAlbForwardingRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultAlbForwardingRuleCols
	}

	columnsMap := map[string]string{
		"ForwardingRuleId":   "ForwardingRuleId",
		"Name":               "Name",
		"Protocol":           "Protocol",
		"ListenerIp":         "ListenerIp",
		"ListenerPort":       "ListenerPort",
		"ServerCertificates": "ServerCertificates",
		"State":              "State",
	}
	var forwardingRuleCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			forwardingRuleCols = append(forwardingRuleCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return forwardingRuleCols
}

func getAlbForwardingRules(forwardingrules resources.ApplicationLoadBalancerForwardingRules) []resources.ApplicationLoadBalancerForwardingRule {
	ss := make([]resources.ApplicationLoadBalancerForwardingRule, 0)
	for _, s := range *forwardingrules.Items {
		ss = append(ss, resources.ApplicationLoadBalancerForwardingRule{ApplicationLoadBalancerForwardingRule: s})
	}
	return ss
}

func getAlbForwardingRulesKVMaps(ss []resources.ApplicationLoadBalancerForwardingRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var forwardingRulePrint AlbForwardingRulePrint
		if idOk, ok := s.GetIdOk(); ok && idOk != nil {
			forwardingRulePrint.ForwardingRuleId = *idOk
		}
		if propertiesOk, ok := s.GetPropertiesOk(); ok && propertiesOk != nil {
			if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
				forwardingRulePrint.Name = *nameOk
			}
			if listenerIpOk, ok := propertiesOk.GetListenerIpOk(); ok && listenerIpOk != nil {
				forwardingRulePrint.ListenerIp = *listenerIpOk
			}
			if listenerPortOk, ok := propertiesOk.GetListenerPortOk(); ok && listenerPortOk != nil {
				forwardingRulePrint.ListenerPort = *listenerPortOk
			}
			if protocolOk, ok := propertiesOk.GetProtocolOk(); ok && protocolOk != nil {
				forwardingRulePrint.Protocol = *protocolOk
			}
			if serverCertificatesOk, ok := propertiesOk.GetServerCertificatesOk(); ok && serverCertificatesOk != nil {
				forwardingRulePrint.ServerCertificates = *serverCertificatesOk
			}
		}
		if metadataOk, ok := s.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				forwardingRulePrint.State = *stateOk
			}
		}
		o := structs.Map(forwardingRulePrint)
		out = append(out, o)
	}
	return out
}

func getAlbForwardingRulesIds(outErr io.Writer, datacenterId, albId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	albSvc := resources.NewApplicationLoadBalancerService(clientSvc.Get(), context.TODO())
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.ApplicationLoadBalancerForwardingRules.GetItemsOk(); ok && items != nil {
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
