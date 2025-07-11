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
	defaultNatGatewayRuleCols = []string{"NatGatewayRuleId", "Name", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "State"}
	allNatGatewayRuleCols     = []string{"NatGatewayRuleId", "Name", "Type", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "TargetPortRangeStart", "TargetPortRangeEnd", "State"}
)

func NatgatewayRuleCmd() *core.Command {
	ctx := context.TODO()
	natgatewayRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r"},
			Short:            "NAT Gateway Rule Operations",
			Long:             "The sub-commands of `ionosctl natgateway rule` allow you to create, list, get, update, delete NAT Gateway Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := natgatewayRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.FlagCols, "", defaultNatGatewayRuleCols, tabheaders.ColsMessage(allNatGatewayRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(natgatewayRuleCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = natgatewayRuleCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNatGatewayRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, natgatewayRuleCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "rule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List NAT Gateway Rules",
		LongDesc:   "Use this command to list NAT Gateway Rules from a specified NAT Gateway.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.NATGatewayRulesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id",
		Example:    listNatGatewayRuleExample,
		PreCmdRun:  PreRunNATGatewayRuleList,
		CmdRun:     RunNatGatewayRuleList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NATGatewayRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NATGatewayRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, natgatewayRuleCmd, core.CommandBuilder{
		Namespace:  "natgateway",
		Resource:   "rule",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a NAT Gateway Rule",
		LongDesc:   "Use this command to get information about a specified NAT Gateway Rule from a NAT Gateway.\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id\n* NAT Gateway Rule Id",
		Example:    getNatGatewayRuleExample,
		PreCmdRun:  PreRunDcNatGatewayRuleIds,
		CmdRun:     RunNatGatewayRuleGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, natgatewayRuleCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NAT Gateway Rule",
		LongDesc: `Use this command to create a NAT Gateway Rule in a specified NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* Public IP
* Source Subnet`,
		Example:    createNatGatewayRuleExample,
		PreCmdRun:  PreRunNatGatewayRuleCreate,
		CmdRun:     RunNatGatewayRuleCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed Rule", "Name of the NAT Gateway Rule")
	create.AddStringFlag(cloudapiv6.FlagProtocol, cloudapiv6.FlagProtocolShort, string(ionoscloud.ALL), "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIpFlag(cloudapiv6.FlagIp, "", nil, "Public IP address of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.FlagSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.FlagTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.FlagPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.FlagPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule creation [seconds]")
	create.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, natgatewayRuleCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a NAT Gateway Rule",
		LongDesc: `Use this command to update a specified NAT Gateway Rule from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id`,
		Example:    updateNatGatewayRuleExample,
		PreCmdRun:  PreRunDcNatGatewayRuleIds,
		CmdRun:     RunNatGatewayRuleUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "Name of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.FlagProtocol, cloudapiv6.FlagProtocolShort, "", "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIpFlag(cloudapiv6.FlagIp, "", nil, "Public IP address of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.FlagSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.FlagTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	update.AddIntFlag(cloudapiv6.FlagPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	update.AddIntFlag(cloudapiv6.FlagPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule update [seconds]")
	update.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, natgatewayRuleCmd, core.CommandBuilder{
		Namespace: "natgateway",
		Resource:  "rule",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a NAT Gateway Rule",
		LongDesc: `Use this command to delete a specified NAT Gateway Rule from a NAT Gateway.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* NAT Gateway Id
* NAT Gateway Rule Id`,
		Example:    deleteNatGatewayRuleExample,
		PreCmdRun:  PreRunDcNatGatewayRuleDelete,
		CmdRun:     RunNatGatewayRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all NAT Gateway Rules.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	return core.WithConfigOverride(natgatewayRuleCmd, "compute", "")
}

func PreRunNATGatewayRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.NATGatewayRulesFilters(), completer.NATGatewayRulesFiltersUsage())
	}
	return nil
}

func PreRunNatGatewayRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagIp, cloudapiv6.FlagSourceSubnet)
}

func PreRunDcNatGatewayRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagRuleId)
}

func PreRunDcNatGatewayRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagRuleId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagNatGatewayId, cloudapiv6.FlagAll},
	)
}

func RunNatGatewayRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	natgatewayRules, resp, err := c.CloudApiV6Services.NatGateways().ListRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.NatGatewayRule, natgatewayRules.NatGatewayRules,
		tabheaders.GetHeaders(allNatGatewayRuleCols, defaultNatGatewayRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayRuleGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"atGatewayRule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))

	ng, resp, err := c.CloudApiV6Services.NatGateways().GetRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NatGatewayRule, ng.NatGatewayRule,
		tabheaders.GetHeaders(allNatGatewayRuleCols, defaultNatGatewayRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayRuleCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	proper := getNewNatGatewayRuleInfo(c)

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))))
	}

	if !proper.HasProtocol() {
		proper.SetProtocol(ionoscloud.NatGatewayRuleProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol))))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol))))
	}

	ng, resp, err := c.CloudApiV6Services.NatGateways().CreateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		resources.NatGatewayRule{
			NatGatewayRule: ionoscloud.NatGatewayRule{
				Properties: &proper.NatGatewayRuleProperties,
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NatGatewayRule, ng.NatGatewayRule,
		tabheaders.GetHeaders(allNatGatewayRuleCols, defaultNatGatewayRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayRuleUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	input := getNewNatGatewayRuleInfo(c)
	ng, resp, err := c.CloudApiV6Services.NatGateways().UpdateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		*input,
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NatGatewayRule, ng.NatGatewayRule,
		tabheaders.GetHeaders(allNatGatewayRuleCols, defaultNatGatewayRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunNatGatewayRuleDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllNatgatewayRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nat gateway rule", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting NatGatewayRule with id: %v...", ruleId))

	resp, err := c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, ruleId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("NAT Gateway Rule successfully deleted"))
	return nil
}

func getNewNatGatewayRuleInfo(c *core.CommandConfig) *resources.NatGatewayRuleProperties {
	input := ionoscloud.NatGatewayRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
		input.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIp)) {
		publicIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIp))
		input.SetPublicIp(publicIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PublicIp set: %v", publicIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol)) {
		protocol := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol)))
		input.SetProtocol(ionoscloud.NatGatewayRuleProtocol(protocol))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Protocol set: %v", protocol))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagSourceSubnet)) {
		sourceSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSourceSubnet))
		input.SetSourceSubnet(sourceSubnet)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SourceSubnet set: %v", sourceSubnet))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagTargetSubnet)) {
		targetSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagTargetSubnet))
		input.SetTargetSubnet(targetSubnet)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", targetSubnet))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeStart)) &&
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeEnd)) {
		inputPortRange := ionoscloud.TargetPortRange{}

		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeStart))
		portRangeStop := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeEnd))

		inputPortRange.SetStart(portRangeStart)
		inputPortRange.SetEnd(portRangeStop)
		input.SetTargetPortRange(inputPortRange)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property TargetPortRang set with start: %v and stop: %v", portRangeStart, portRangeStop))
	}

	return &resources.NatGatewayRuleProperties{
		NatGatewayRuleProperties: input,
	}
}

func DeleteAllNatgatewayRules(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNatGatewayId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("NatGateway ID: %v", natGatewayId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting NatGateway Rules..."))

	natGatewayRules, resp, err := c.CloudApiV6Services.NatGateways().ListRules(dcId, natGatewayId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	natGatewayRuleItems, ok := natGatewayRules.GetItemsOk()
	if !ok || natGatewayRuleItems == nil {
		return fmt.Errorf("could not get items of NAT Gateway Rules")
	}

	if len(*natGatewayRuleItems) <= 0 {
		return fmt.Errorf("no NAT Gateway Rules found")
	}

	var multiErr error
	for _, natGateway := range *natGatewayRuleItems {
		id := natGateway.GetId()
		name := natGateway.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the NatGatewayRule with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
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
