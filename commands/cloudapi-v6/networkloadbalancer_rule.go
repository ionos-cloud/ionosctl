package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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

func NetworkloadbalancerRuleCmd() *core.Command {
	ctx := context.TODO()
	nlbRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Network Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl nlb rule` allow you to create, list, get, update, delete Network Load Balancer Forwarding Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultForwardingRuleCols, printer.ColsMessage(allForwardingRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nlbRuleCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nlbRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "rule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancer Forwarding Rules",
		LongDesc:   "Use this command to list Network Load Balancer Forwarding Rules from a specified Network Load Balancer.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.NlbRulesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id",
		Example:    listNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleList,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "rule",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Network Load Balancer Forwarding Rule",
		LongDesc:   "Use this command to get information about a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id\n* Forwarding Rule Id",
		Example:    getNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a Network Load Balancer Forwarding Rule in a specified Network Load Balancer. You can also set Health Check Settings for Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Listener Ip
* Listener Port`,
		Example:    createNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunNetworkLoadBalancerForwardingRuleCreate,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name for the Forwarding Rule")
	create.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Listening IP", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgListenerPort, "", "", "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	create.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	create.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	create.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	create.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	create.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to update a specified Network Load Balancer Forwarding Rule from a Network Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example:    updateNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Forwarding Rule")
	update.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Listening IP", core.RequiredFlagOption())
	update.AddStringFlag(cloudapiv6.ArgListenerPort, "", "", "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	update.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	update.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	update.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	update.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to delete a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example:    deleteNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleDelete,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Network Load Balancer Forwarding Rule.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	nlbRuleCmd.AddCommand(NlbRuleTargetCmd())

	return nlbRuleCmd
}

func PreRunNetworkLoadBalancerRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.NlbRulesFilters(), completer.NlbRulesFiltersUsage())
	}
	return nil
}

func PreRunNetworkLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgListenerIp, cloudapiv6.ArgListenerPort)
}

func PreRunDcNetworkLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId)
}

func PreRunDcNetworkLoadBalancerForwardingRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func RunNetworkLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(nil, c, getForwardingRules(nlbForwardingRules)))
}

func RunNetworkLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("NetworkLoadBalancerForwardingRule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(nil, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	proper := getForwardingRulePropertiesSet(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(string(ionoscloud.TCP))
	}
	if !proper.HasAlgorithm() {
		proper.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	}
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		proper.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		resources.NetworkLoadBalancerForwardingRule{
			NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
				Properties: &proper.NetworkLoadBalancerForwardingRuleProperties,
			},
		},
		queryParams,
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
	return c.Printer.Print(getForwardingRulePrint(resp, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	input := getForwardingRulePropertiesSet(c)
	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		input.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		input,
		queryParams,
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
	return c.Printer.Print(getForwardingRulePrint(resp, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNetworkLoadBalancerForwardingRules(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete network load balancer forwarding rule"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting NetworkLoadBalancerForwardingRule with id: %v...", ruleId)
		resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, ruleId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getForwardingRulePrint(resp, c, nil))
	}
}

func getForwardingRulePropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleProperties {
	input := ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)) {
		algorithm := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
		input.SetAlgorithm(algorithm)
		c.Printer.Verbose("Property Algorithm set: %v", algorithm)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)) {
		listenerIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp))
		input.SetListenerIp(listenerIp)
		c.Printer.Verbose("Property ListenerIp set: %v", listenerIp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)) {
		listenerPort := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort))
		input.SetListenerPort(listenerPort)
		c.Printer.Verbose("Property ListenerPort set: %v", listenerPort)
	}
	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: input,
	}
}

func getForwardingRuleHealthCheckPropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleHealthCheck {
	inputHealth := ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealth.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)) {
		inputHealth.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)) {
		inputHealth.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)) {
		inputHealth.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
	}
	return &resources.NetworkLoadBalancerForwardingRuleHealthCheck{
		NetworkLoadBalancerForwardingRuleHealthCheck: inputHealth,
	}
}

func DeleteAllNetworkLoadBalancerForwardingRules(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("NetworkLoadBalancer ID: %v", loadBalancerId)
	c.Printer.Verbose("Getting NetworkLoadBalancerForwardingRules...")
	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(dcId, loadBalancerId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if nlbForwardingRulesItems, ok := nlbForwardingRules.GetItemsOk(); ok && nlbForwardingRulesItems != nil {
		if len(*nlbForwardingRulesItems) > 0 {
			_ = c.Printer.Warn("NetworkLoadBalancerForwardingRules to be deleted:")
			for _, nlbForwardingRule := range *nlbForwardingRulesItems {
				delIdAndName := ""
				if id, ok := nlbForwardingRule.GetIdOk(); ok && id != nil {
					delIdAndName += "NetworkLoadBalancerForwardingRule Id: " + *id
				}
				if properties, ok := nlbForwardingRule.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += " NetworkLoadBalancerForwardingRule Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Backup Units"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the NetworkLoadBalancerForwardingRules...")
			var multiErr error
			for _, nlbForwardingRule := range *nlbForwardingRulesItems {
				if id, ok := nlbForwardingRule.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting NetworkLoadBalancerForwardingRule with id: %v...", *id)
					resp, err = c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
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
			return errors.New("no NetworkLoadBalancerForwardingRules found")
		}
	} else {
		return errors.New("could not get items of NetworkLoadBalancerForwardingRules")
	}
}

// Output Printing

var (
	defaultForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State"}
	allForwardingRuleCols     = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State",
		"ClientTimeout", "ConnectTimeout", "TargetTimeout", "Retries"}
)

type ForwardingRulePrint struct {
	ForwardingRuleId string `json:"ForwardingRuleId,omitempty"`
	Name             string `json:"Name,omitempty"`
	Algorithm        string `json:"Algorithm,omitempty"`
	Protocol         string `json:"Protocol,omitempty"`
	ListenerIp       string `json:"ListenerIp,omitempty"`
	ListenerPort     int32  `json:"ListenerPort,omitempty"`
	ClientTimeout    string `json:"ClientTimeout,omitempty"`
	ConnectTimeout   string `json:"ConnectTimeout,omitempty"`
	TargetTimeout    string `json:"TargetTimeout,omitempty"`
	Retries          int32  `json:"Retries,omitempty"`
	State            string `json:"State,omitempty"`
}

func getForwardingRulePrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NetworkLoadBalancerForwardingRule) printer.Result {
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
			r.KeyValue = getForwardingRulesKVMaps(ss)
			r.Columns = getForwardingRulesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getForwardingRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultForwardingRuleCols
	}

	columnsMap := map[string]string{
		"ForwardingRuleId": "ForwardingRuleId",
		"Name":             "Name",
		"Algorithm":        "Algorithm",
		"Protocol":         "Protocol",
		"ListenerIp":       "ListenerIp",
		"ListenerPort":     "ListenerPort",
		"ClientTimeout":    "ClientTimeout",
		"ConnectTimeout":   "ConnectTimeout",
		"TargetTimeout":    "TargetTimeout",
		"Retries":          "Retries",
		"State":            "State",
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

func getForwardingRules(forwardingrules resources.NetworkLoadBalancerForwardingRules) []resources.NetworkLoadBalancerForwardingRule {
	nlbRuleObjs := make([]resources.NetworkLoadBalancerForwardingRule, 0)
	if items, ok := forwardingrules.GetItemsOk(); ok && items != nil {
		for _, forwardingRule := range *items {
			nlbRuleObjs = append(nlbRuleObjs, resources.NetworkLoadBalancerForwardingRule{NetworkLoadBalancerForwardingRule: forwardingRule})
		}
	}
	return nlbRuleObjs
}

func getForwardingRulesKVMaps(ss []resources.NetworkLoadBalancerForwardingRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var forwardingRulePrint ForwardingRulePrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			forwardingRulePrint.ForwardingRuleId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				forwardingRulePrint.Name = *name
			}
			if listenerIp, ok := properties.GetListenerIpOk(); ok && listenerIp != nil {
				forwardingRulePrint.ListenerIp = *listenerIp
			}
			if listenerPort, ok := properties.GetListenerPortOk(); ok && listenerPort != nil {
				forwardingRulePrint.ListenerPort = *listenerPort
			}
			if protocol, ok := properties.GetProtocolOk(); ok && protocol != nil {
				forwardingRulePrint.Protocol = *protocol
			}
			if algorithm, ok := properties.GetAlgorithmOk(); ok && algorithm != nil {
				forwardingRulePrint.Algorithm = *algorithm
			}
			if health, ok := properties.GetHealthCheckOk(); ok && health != nil {
				if clientTimeout, ok := health.GetClientTimeoutOk(); ok && clientTimeout != nil {
					forwardingRulePrint.ClientTimeout = fmt.Sprintf("%vms", *clientTimeout)
				}
				if connectTimeout, ok := health.GetConnectTimeoutOk(); ok && connectTimeout != nil {
					forwardingRulePrint.ConnectTimeout = fmt.Sprintf("%vms", *connectTimeout)
				}
				if targetTimeout, ok := health.GetTargetTimeoutOk(); ok && targetTimeout != nil {
					forwardingRulePrint.TargetTimeout = fmt.Sprintf("%vms", *targetTimeout)
				}
				if retries, ok := health.GetRetriesOk(); ok && retries != nil {
					forwardingRulePrint.Retries = *retries
				}
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				forwardingRulePrint.State = *state
			}
		}
		o := structs.Map(forwardingRulePrint)
		out = append(out, o)
	}
	return out
}
