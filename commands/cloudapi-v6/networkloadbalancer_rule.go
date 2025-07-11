package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State"}
	allForwardingRuleCols     = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State",
		"ClientTimeout", "ConnectTimeout", "TargetTimeout", "Retries"}
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
	globalFlags.StringSliceP(constants.FlagCols, "", defaultForwardingRuleCols, tabheaders.ColsMessage(allForwardingRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(nlbRuleCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = nlbRuleCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NlbRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})

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
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

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
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed Forwarding Rule", "The name for the Forwarding Rule")
	create.AddIpFlag(cloudapiv6.FlagListenerIp, "", nil, "Listening IP", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.FlagListenerPort, "", "", "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	create.AddIntFlag(cloudapiv6.FlagRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	create.AddIntFlag(cloudapiv6.FlagClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	create.AddIntFlag(cloudapiv6.FlagConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	create.AddIntFlag(cloudapiv6.FlagTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	create.AddStringFlag(cloudapiv6.FlagAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

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
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "The name for the Forwarding Rule")
	update.AddIpFlag(cloudapiv6.FlagListenerIp, "", nil, "Listening IP", core.RequiredFlagOption())
	update.AddStringFlag(cloudapiv6.FlagListenerPort, "", "", "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	update.AddStringFlag(cloudapiv6.FlagAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.FlagRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	update.AddIntFlag(cloudapiv6.FlagClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	update.AddIntFlag(cloudapiv6.FlagConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	update.AddIntFlag(cloudapiv6.FlagTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all Network Load Balancer Forwarding Rule.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	nlbRuleCmd.AddCommand(NlbRuleTargetCmd())

	return core.WithConfigOverride(nlbRuleCmd, "compute", "")
}

func PreRunNetworkLoadBalancerRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNetworkLoadBalancerId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.NlbRulesFilters(), completer.NlbRulesFiltersUsage())
	}
	return nil
}

func PreRunNetworkLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNetworkLoadBalancerId, cloudapiv6.FlagListenerIp, cloudapiv6.FlagListenerPort)
}

func PreRunDcNetworkLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNetworkLoadBalancerId, cloudapiv6.FlagRuleId)
}

func PreRunDcNetworkLoadBalancerForwardingRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNetworkLoadBalancerId, cloudapiv6.FlagRuleId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNetworkLoadBalancerId, cloudapiv6.FlagAll},
	)
}

func RunNetworkLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.NetworkLoadBalancerRule,
		nlbForwardingRules.NetworkLoadBalancerForwardingRules, tabheaders.GetHeaders(allForwardingRuleCols, defaultForwardingRuleCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNetworkLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Network Load Balancer Forwarding Rule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NetworkLoadBalancerRule,
		ng.NetworkLoadBalancerForwardingRule, tabheaders.GetHeaders(allForwardingRuleCols, defaultForwardingRuleCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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
		proper.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagAlgorithm)))
	}

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
	}

	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		proper.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId)),
		resources.NetworkLoadBalancerForwardingRule{
			NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
				Properties: &proper.NetworkLoadBalancerForwardingRuleProperties,
			},
		},
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NetworkLoadBalancerRule,
		ng.NetworkLoadBalancerForwardingRule, tabheaders.GetHeaders(allForwardingRuleCols, defaultForwardingRuleCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		input,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NetworkLoadBalancerRule,
		ng.NetworkLoadBalancerForwardingRule, tabheaders.GetHeaders(allForwardingRuleCols, defaultForwardingRuleCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNetworkLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllNetworkLoadBalancerForwardingRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete network load balancer forwarding rule", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Network Load Balancer Forwarding Rule with id: %v...", ruleId))

	resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, ruleId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Network Load Balancer Forwarding Rule successfully deleted"))
	return nil
}

func getForwardingRulePropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleProperties {
	input := ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagAlgorithm)) {
		algorithm := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagAlgorithm)))
		input.SetAlgorithm(algorithm)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Algorithm set: %v", algorithm))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagListenerIp)) {
		listenerIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagListenerIp))
		input.SetListenerIp(listenerIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property ListenerIp set: %v", listenerIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagListenerPort)) {
		listenerPort := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagListenerPort))
		input.SetListenerPort(listenerPort)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property ListenerPort set: %v", listenerPort))
	}

	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: input,
	}
}

func getForwardingRuleHealthCheckPropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleHealthCheck {
	inputHealth := ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagRetries)) {
		inputHealth.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagRetries)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagClientTimeout)) {
		inputHealth.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagClientTimeout)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagConnectionTimeout)) {
		inputHealth.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagConnectionTimeout)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagTargetTimeout)) {
		inputHealth.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagTargetTimeout)))
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
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNetworkLoadBalancerId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Network Load Balancer ID: %v", loadBalancerId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Network Load Balancer Forwarding Rules..."))

	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(dcId, loadBalancerId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	nlbForwardingRulesItems, ok := nlbForwardingRules.GetItemsOk()
	if !ok || nlbForwardingRulesItems == nil {
		return fmt.Errorf("could not get items of Network Load Balancer Forwarding Rules")
	}

	if len(*nlbForwardingRulesItems) <= 0 {
		return fmt.Errorf("no Network Load Balancer Forwarding Rules found")
	}

	var multiErr error
	for _, nlbForwardingRule := range *nlbForwardingRulesItems {
		name := nlbForwardingRule.GetProperties().Name
		id := nlbForwardingRule.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Network Load Balancer Forwarding Rule with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
