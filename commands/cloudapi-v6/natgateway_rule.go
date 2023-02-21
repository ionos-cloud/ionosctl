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
	"github.com/ionos-cloud/ionosctl/pkg/constants"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNatGatewayRuleCols, printer.ColsMessage(allNatGatewayRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(natgatewayRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = natgatewayRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NATGatewayRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NATGatewayRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Rule", "Name of the NAT Gateway Rule")
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, string(ionoscloud.ALL), "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIpFlag(cloudapiv6.ArgIp, "", nil, "Public IP address of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv6.ArgTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	create.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "", "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIpFlag(cloudapiv6.ArgIp, "", nil, "Public IP address of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.ArgSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule")
	update.AddStringFlag(cloudapiv6.ArgTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	update.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	update.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgNatGatewayId, "", "", cloudapiv6.NatGatewayId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.RuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for NAT Gateway Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all NAT Gateway Rules.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return natgatewayRuleCmd
}

func PreRunNATGatewayRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.NATGatewayRulesFilters(), completer.NATGatewayRulesFiltersUsage())
	}
	return nil
}

func PreRunNatGatewayRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIp, cloudapiv6.ArgSourceSubnet)
}

func PreRunDcNatGatewayRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgRuleId)
}

func PreRunDcNatGatewayRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	natgatewayRules, resp, err := c.CloudApiV6Services.NatGateways().ListRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(nil, c, getNatGatewayRules(natgatewayRules)))
}

func RunNatGatewayRuleGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("atGatewayRule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	ng, resp, err := c.CloudApiV6Services.NatGateways().GetRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(nil, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	proper := getNewNatGatewayRuleInfo(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if !proper.HasProtocol() {
		proper.SetProtocol(ionoscloud.NatGatewayRuleProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))
		c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}
	ng, resp, err := c.CloudApiV6Services.NatGateways().CreateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		resources.NatGatewayRule{
			NatGatewayRule: ionoscloud.NatGatewayRule{
				Properties: &proper.NatGatewayRuleProperties,
			},
		},
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(resp, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	input := getNewNatGatewayRuleInfo(c)
	ng, resp, err := c.CloudApiV6Services.NatGateways().UpdateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		*input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(resp, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNatgatewayRules(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nat gateway rule"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting NatGatewayRule with id: %v...", ruleId)
		resp, err := c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, ruleId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getNatGatewayRulePrint(resp, c, nil))
	}
}

func getNewNatGatewayRuleInfo(c *core.CommandConfig) *resources.NatGatewayRuleProperties {
	input := ionoscloud.NatGatewayRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIp)) {
		publicIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
		input.SetPublicIp(publicIp)
		c.Printer.Verbose("Property PublicIp set: %v", publicIp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		protocol := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		input.SetProtocol(ionoscloud.NatGatewayRuleProtocol(protocol))
		c.Printer.Verbose("Property Protocol set: %v", protocol)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceSubnet)) {
		sourceSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceSubnet))
		input.SetSourceSubnet(sourceSubnet)
		c.Printer.Verbose("Property SourceSubnet set: %v", sourceSubnet)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetSubnet)) {
		targetSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetSubnet))
		input.SetTargetSubnet(targetSubnet)
		c.Printer.Verbose("Property Name set: %v", targetSubnet)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart)) &&
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd)) {
		inputPortRange := ionoscloud.TargetPortRange{}
		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart))
		portRangeStop := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd))
		inputPortRange.SetStart(portRangeStart)
		inputPortRange.SetEnd(portRangeStop)
		input.SetTargetPortRange(inputPortRange)
		c.Printer.Verbose("Property TargetPortRang set with start: %v and stop: %v", portRangeStart, portRangeStop)
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
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("NatGateway ID: %v", natGatewayId)
	c.Printer.Verbose("Getting NatGateway Rules...")
	natGatewayRules, resp, err := c.CloudApiV6Services.NatGateways().ListRules(dcId, natGatewayId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if natGatewayRuleItems, ok := natGatewayRules.GetItemsOk(); ok && natGatewayRuleItems != nil {
		if len(*natGatewayRuleItems) > 0 {
			_ = c.Printer.Warn("NatGatewayRules to be deleted:")
			for _, natGateway := range *natGatewayRuleItems {
				delIdAndName := ""
				if id, ok := natGateway.GetIdOk(); ok && id != nil {
					delIdAndName += "NatGatewayRule Id: " + *id
				}
				if properties, ok := natGateway.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += " NatGatewayRule Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the NatGatewayRules"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the NatGatewayRules...")
			var multiErr error
			for _, natGateway := range *natGatewayRuleItems {
				if id, ok := natGateway.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting NatGatewayRule with id: %v...", *id)
					resp, err = c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no NatGatewayRules found")
		}
	} else {
		return errors.New("could not get items of NatGatewayRules")
	}
}

// Output Printing

var (
	defaultNatGatewayRuleCols = []string{"NatGatewayRuleId", "Name", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "State"}
	allNatGatewayRuleCols     = []string{"NatGatewayRuleId", "Name", "Type", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "TargetPortRangeStart", "TargetPortRangeEnd", "State"}
)

type NatGatewayRulePrint struct {
	NatGatewayRuleId     string `json:"NatGatewayRuleId,omitempty"`
	Name                 string `json:"Name,omitempty"`
	Type                 string `json:"Type,omitempty"`
	Protocol             string `json:"Protocol,omitempty"`
	SourceSubnet         string `json:"SourceSubnet,omitempty"`
	PublicIp             string `json:"PublicIp,omitempty"`
	TargetSubnet         string `json:"TargetSubnet,omitempty"`
	TargetPortRangeStart int32  `json:"TargetPortRangeStart,omitempty"`
	TargetPortRangeEnd   int32  `json:"TargetPortRangeEnd,omitempty"`
	State                string `json:"State,omitempty"`
}

func getNatGatewayRulePrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NatGatewayRule) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getNatGatewayRulesKVMaps(ss)
			r.Columns = printer.GetHeaders(allNatGatewayRuleCols, defaultNatGatewayRuleCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols)))
		}
	}
	return r
}

func getNatGatewayRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultNatGatewayRuleCols
	}

	columnsMap := map[string]string{
		"NatGatewayId":         "NatGatewayId",
		"Name":                 "Name",
		"PublicIp":             "PublicIp",
		"Type":                 "Type",
		"Protocol":             "Protocol",
		"SourceSubnet":         "SourceSubnet",
		"TargetSubnet":         "TargetSubnet",
		"TargetPortRangeStart": "TargetPortRangeStart",
		"TargetPortRangeEnd":   "TargetPortRangeEnd",
		"State":                "State",
	}
	var natgatewayRuleCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			natgatewayRuleCols = append(natgatewayRuleCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return natgatewayRuleCols
}

func getNatGatewayRules(natgatewayRules resources.NatGatewayRules) []resources.NatGatewayRule {
	ruleObjs := make([]resources.NatGatewayRule, 0)
	if items, ok := natgatewayRules.GetItemsOk(); ok && items != nil {
		for _, natGatewayRule := range *items {
			ruleObjs = append(ruleObjs, resources.NatGatewayRule{NatGatewayRule: natGatewayRule})
		}
	}
	return ruleObjs
}

func getNatGatewayRulesKVMaps(ss []resources.NatGatewayRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var natgatewayRulePrint NatGatewayRulePrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			natgatewayRulePrint.NatGatewayRuleId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				natgatewayRulePrint.Name = *name
			}
			if t, ok := properties.GetTypeOk(); ok && t != nil {
				natgatewayRulePrint.Type = string(*t)
			}
			if protocol, ok := properties.GetProtocolOk(); ok && protocol != nil {
				natgatewayRulePrint.Protocol = string(*protocol)
			}
			if ip, ok := properties.GetPublicIpOk(); ok && ip != nil {
				natgatewayRulePrint.PublicIp = *ip
			}
			if ssubnet, ok := properties.GetSourceSubnetOk(); ok && ssubnet != nil {
				natgatewayRulePrint.SourceSubnet = *ssubnet
			}
			if tsubnet, ok := properties.GetTargetSubnetOk(); ok && tsubnet != nil {
				natgatewayRulePrint.TargetSubnet = *tsubnet
			}
			if portRange, ok := properties.GetTargetPortRangeOk(); ok && portRange != nil {
				if portRangeStart, ok := portRange.GetStartOk(); ok && portRangeStart != nil {
					natgatewayRulePrint.TargetPortRangeStart = *portRangeStart
				}
				if portRangeEnd, ok := portRange.GetEndOk(); ok && portRangeEnd != nil {
					natgatewayRulePrint.TargetPortRangeEnd = *portRangeEnd
				}
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				natgatewayRulePrint.State = *state
			}
		}
		o := structs.Map(natgatewayRulePrint)
		out = append(out, o)
	}
	return out
}
