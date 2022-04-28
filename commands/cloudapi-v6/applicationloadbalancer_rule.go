package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"go.uber.org/multierr"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
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
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
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
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
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
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name of the Application Load Balancer forwarding rule.")
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Balancing protocol.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgListenerIp, "", "", "Listening (inbound) IP.", core.RequiredFlagOption())
	create.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening (inbound) port number; valid range is 1 to 65535.", core.RequiredFlagOption())
	create.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 50, "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).")
	create.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
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
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the Application Load Balancer forwarding rule.")
	update.AddStringFlag(cloudapiv6.ArgListenerIp, "", "", "Listening (inbound) IP.")
	update.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening (inbound) port number; valid range is 1 to 65535.")
	update.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 50, "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).")
	update.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, cloudapiv6.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultAlbForwardingRuleCols, printer.ColsMessage(defaultAlbForwardingRuleCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunApplicationLoadBalancerForwardingRuleDelete,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Forwarding Rules")
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule deletion [seconds]")

	albRuleCmd.AddCommand(AlbRuleHttpRuleCmd())

	return albRuleCmd
}

func PreRunApplicationLoadBalancerForwardingRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func PreRunApplicationLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgListenerIp, cloudapiv6.ArgListenerPort)
}

func PreRunDcApplicationLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId)
}

func RunApplicationLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	c.Printer.Verbose("Getting ForwardingRules")
	albForwardingRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(nil, c, getAlbForwardingRules(albForwardingRules)))
}

func RunApplicationLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("Getting ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(nil, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	proper := getAlbForwardingRulePropertiesSet(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	c.Printer.Verbose("Creating ForwardingRule")
	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.ApplicationLoadBalancerForwardingRule{
			ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
				Properties: &proper.ApplicationLoadBalancerForwardingRuleProperties,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	input := getAlbForwardingRulePropertiesSet(c)
	c.Printer.Verbose("Updating ForwardingRule")
	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, []resources.ApplicationLoadBalancerForwardingRule{*applicationLoadBalancerForwardingRule}))
}

func RunApplicationLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	var (
		resp *resources.Response
		err  error
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		err = DeleteAllApplicationLoadBalancerForwardingRule(c)
		if err != nil {
			return err
		}
	} else {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete application load balancer forwarding rule"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getAlbForwardingRulePrint(resp, c, nil))
}

func DeleteAllApplicationLoadBalancerForwardingRule(c *core.CommandConfig) error {
	_ = c.Printer.Print("Getting Application Load Balancer Forwarding Rules...")
	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.ListQueryParams{},
	)
	if err != nil {
		return err
	}
	if albRuleItems, ok := applicationLoadBalancerRules.GetItemsOk(); ok && albRuleItems != nil {
		if len(*albRuleItems) > 0 {
			for _, fr := range *albRuleItems {
				toPrint := ""
				if id, ok := fr.GetIdOk(); ok && id != nil {
					toPrint += "Forwarding Rule Id: " + *id
				}
				if properties, ok := fr.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += "Forwarding Rule Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Forwarding Rules"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Application Load Balancer Forwarding Rules...")
			var multiErr error
			for _, fr := range *albRuleItems {
				if id, ok := fr.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Application Load Balancer Forwarding Rule with id: %v...", *id)
					resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteForwardingRule(
						viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
						viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), *id)
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
			return errors.New("no Target Groups found")
		}
	} else {
		return errors.New("could not get items of Target Groups")
	}
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
		c.Printer.Verbose("Property ListenerPort set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)) {
		input.SetServerCertificates(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
		c.Printer.Verbose("Property ServerCertificates set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)) {
		input.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
		c.Printer.Verbose("Property Client Timeout set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
	}
	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: input,
	}
}

// Output Printing

var defaultAlbForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Protocol", "ListenerIp", "ListenerPort", "ServerCertificates", "State"}
var allAlbForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Protocol", "ListenerIp", "ListenerPort", "ClientTimeout", "ServerCertificates", "State"}

type AlbForwardingRulePrint struct {
	ForwardingRuleId   string   `json:"ForwardingRuleId,omitempty"`
	Name               string   `json:"Name,omitempty"`
	Protocol           string   `json:"Protocol,omitempty"`
	ListenerIp         string   `json:"ListenerIp,omitempty"`
	ListenerPort       int32    `json:"ListenerPort,omitempty"`
	ClientTimeout      int32    `json:"ClientTimeout,omitempty"`
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
		"ClientTimeout":      "ClientTimeout",
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
			if clientTimeoutOk, ok := propertiesOk.GetClientTimeoutOk(); ok && clientTimeoutOk != nil {
				forwardingRulePrint.ClientTimeout = *clientTimeoutOk
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
