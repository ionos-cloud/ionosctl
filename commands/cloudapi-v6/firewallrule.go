package commands

import (
	"context"
	"errors"
	"fmt"
	"net"
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
	defaultFirewallRuleCols = []string{"FirewallRuleId", "Name", "Protocol", "PortRangeStart", "PortRangeEnd", "Direction", "IPVersion", "State"}
	allFirewallRuleCols     = []string{"FirewallRuleId", "Name", "Protocol", "SourceMac", "SourceIP", "DestinationIP", "PortRangeStart", "PortRangeEnd",
		"IcmpCode", "IcmpType", "Direction", "IPVersion", "State"}
)

func FirewallruleCmd() *core.Command {
	ctx := context.TODO()
	firewallRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "firewallrule",
			Aliases:          []string{"f", "fr", "firewall"},
			Short:            "Firewall Rule Operations",
			Long:             "The sub-commands of `ionosctl firewallrule` allow you to create, list, get, update, delete Firewall Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := firewallRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.FlagCols, "", defaultFirewallRuleCols, tabheaders.ColsMessage(allFirewallRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(firewallRuleCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allFirewallRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, firewallRuleCmd, core.CommandBuilder{
		Namespace:  "firewallrule",
		Resource:   "firewallrule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Firewall Rules",
		LongDesc:   "Use this command to get a list of Firewall Rules from a specified NIC from a Server.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.FirewallRulesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Nic Id",
		Example:    listFirewallRuleExample,
		PreCmdRun:  PreRunFirewallRuleList,
		CmdRun:     RunFirewallRuleList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, "", cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)
	list.AddStringFlag(cloudapiv6.FlagOrderBy, "", "", cloudapiv6.FlagOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.FlagFilters, cloudapiv6.FlagFiltersShort, []string{""}, cloudapiv6.FlagFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, firewallRuleCmd, core.CommandBuilder{
		Namespace:  "firewallrule",
		Resource:   "firewallrule",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Firewall Rule",
		LongDesc:   "Use this command to retrieve information of a specified Firewall Rule.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Nic Id\n* FirewallRule Id",
		Example:    getFirewallRuleExample,
		PreCmdRun:  PreRunDcServerNicFRuleIds,
		CmdRun:     RunFirewallRuleGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.FlagFirewallRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagServerId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.FlagNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.FlagDepth, "", cloudapiv6.DefaultGetDepth, cloudapiv6.FlagDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, firewallRuleCmd, core.CommandBuilder{
		Namespace: "firewallrule",
		Resource:  "firewallrule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Firewall Rule",
		LongDesc: `Use this command to create/add a new Firewall Rule to the specified NIC. All Firewall Rules must be associated with a NIC.

NOTE: the Firewall Rule Protocol can only be set when creating a new Firewall Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Protocol`,
		Example:    createFirewallRuleExample,
		PreCmdRun:  PreRunDcServerNicIdsFRuleProtocol,
		CmdRun:     RunFirewallRuleCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "Unnamed Rule", "The name for the Firewall Rule")
	create.AddStringFlag(cloudapiv6.FlagProtocol, "", "", "The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"TCP", "UDP", "ICMP", "ANY"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.FlagSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses")
	create.AddIpFlag(cloudapiv6.FlagSourceIp, "", nil, "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs")
	create.AddIpFlag(cloudapiv6.FlagDestinationIp, cloudapiv6.FlagDestinationIpShort, nil, "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs. WARNING: This short-hand flag `-D` is deprecated.")
	create.AddIntFlag(cloudapiv6.FlagIcmpType, "", 0, "Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types")
	create.AddIntFlag(cloudapiv6.FlagIcmpCode, "", 0, "Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes")
	create.AddIntFlag(cloudapiv6.FlagPortRangeStart, "", 1, "Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	create.AddIntFlag(cloudapiv6.FlagPortRangeEnd, "", 1, "Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	create.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "INGRESS", "The type/direction of Firewall Rule")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule creation to be executed")
	create.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule creation [seconds]")
	create.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.FlagNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(create.NS, cloudapiv6.FlagServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(cloudapiv6.FlagDepth, "", cloudapiv6.DefaultCreateDepth, cloudapiv6.FlagDepthDescription)
	create.AddSetFlag(cloudapiv6.FlagIPVersion, "", "IPv4", []string{"IPv4", "IPv6"}, "The IP version for the Firewall Rule")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, firewallRuleCmd, core.CommandBuilder{
		Namespace: "firewallrule",
		Resource:  "firewallrule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a FirewallRule",
		LongDesc: `Use this command to update a specified Firewall Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id`,
		Example:    updateFirewallRuleExample,
		PreCmdRun:  PreRunDcServerNicFRuleIds,
		CmdRun:     RunFirewallRuleUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "The name for the Firewall Rule")
	update.AddStringFlag(cloudapiv6.FlagSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Not setting option allows all source MAC addresses")
	update.AddIpFlag(cloudapiv6.FlagSourceIp, "", nil, "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs")
	update.AddIpFlag(cloudapiv6.FlagDestinationIp, cloudapiv6.FlagDestinationIpShort, nil, "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs. WARNING: This short-hand flag `-D` is deprecated.")
	update.AddIntFlag(cloudapiv6.FlagIcmpType, "", 0, "Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types")
	update.AddIntFlag(cloudapiv6.FlagIcmpCode, "", 0, "Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes")
	update.AddIntFlag(cloudapiv6.FlagPortRangeStart, "", 1, "Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	update.AddIntFlag(cloudapiv6.FlagPortRangeEnd, "", 1, "Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	update.AddStringFlag(cloudapiv6.FlagDirection, cloudapiv6.FlagDirectionShort, "", "The type/direction of Firewall Rule")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagFirewallRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagServerId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule update to be executed")
	update.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule update [seconds]")
	update.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.FlagNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.FlagServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(cloudapiv6.FlagDepth, "", cloudapiv6.DefaultUpdateDepth, cloudapiv6.FlagDepthDescription)
	update.AddSetFlag(cloudapiv6.FlagIPVersion, "", "IPv4", []string{"IPv4", "IPv6"}, "The IP version for the Firewall Rule")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, firewallRuleCmd, core.CommandBuilder{
		Namespace: "firewallrule",
		Resource:  "firewallrule",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a FirewallRule",
		LongDesc: `Use this command to delete a specified Firewall Rule from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id`,
		Example:    deleteFirewallRuleExample,
		PreCmdRun:  PreRunFirewallDelete,
		CmdRun:     RunFirewallRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.FlagFirewallRuleId, cloudapiv6.FlagIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagServerId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.FlagNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Delete all the Firewalls.")
	deleteCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.FlagDepth, "", cloudapiv6.DefaultDeleteDepth, cloudapiv6.FlagDepthDescription)

	return firewallRuleCmd
}

func PreRunFirewallRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagFilters)) {
		return query.ValidateFilters(c, completer.FirewallRulesFilters(), completer.FirewallRulesFiltersUsage())
	}
	return nil
}

func PreRunDcServerNicIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId)
}

func PreRunFirewallDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId, cloudapiv6.FlagFirewallRuleId},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId, cloudapiv6.FlagAll},
	)
}

func PreRunDcServerNicIdsFRuleProtocol(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId, cloudapiv6.FlagProtocol)
}

func PreRunDcServerNicFRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId, cloudapiv6.FlagServerId, cloudapiv6.FlagNicId, cloudapiv6.FlagFirewallRuleId)
}

func RunFirewallRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	firewallRules, resp, err := c.CloudApiV6Services.FirewallRules().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.FirewallRule, firewallRules.FirewallRules,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunFirewallRuleGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting Firewall Rule with id: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFirewallRuleId))))

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFirewallRuleId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunFirewallRuleCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	properties := getFirewallRulePropertiesSet(c)

	if !properties.HasName() {
		properties.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
	}

	if !properties.HasType() {
		properties.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDirection)))
	}

	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId)),
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunFirewallRuleUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFirewallRuleId)),
		getFirewallRulePropertiesSet(c),
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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunFirewallRuleDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId))
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagFirewallRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		if err := DeleteAllFirewallRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete Firewall Rule with id: %v...", fruleId), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err := c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, fruleId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Firewall Rule successfully deleted"))

	return nil

}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *core.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))
		properties.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol)) {
		protocol := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagProtocol))
		properties.SetProtocol(protocol)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Protocol set: %v", protocol))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagSourceIp)) {
		sourceIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSourceIp))
		properties.SetSourceIp(sourceIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SourceIp set: %v", sourceIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagSourceMac)) {
		sourceMac := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSourceMac))
		properties.SetSourceMac(sourceMac)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SourceMac set: %v", sourceMac))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDestinationIp)) {
		targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDestinationIp))
		properties.SetTargetIp(targetIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property TargetIp/DestinationIp set: %v", targetIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIcmpCode)) {
		icmpCode := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagIcmpCode))
		properties.SetIcmpCode(icmpCode)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IcmpCode set: %v", icmpCode))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIcmpType)) {
		icmpType := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagIcmpType))
		properties.SetIcmpType(icmpType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IcmpType set: %v", icmpType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeStart)) {
		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeStart))
		properties.SetPortRangeStart(portRangeStart)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PortRangeStart set: %v", portRangeStart))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeEnd)) {
		portRangeEnd := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagPortRangeEnd))
		properties.SetPortRangeEnd(portRangeEnd)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PortRangeEnd set: %v", portRangeEnd))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDirection)) {
		firewallruleType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDirection))
		properties.SetType(strings.ToUpper(firewallruleType))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Type/Direction set: %v", firewallruleType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		ipVersion := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion))
		properties.SetIpVersion(ipVersion)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IP Version set: %v", ipVersion))
	}

	return properties
}

func DeleteAllFirewallRules(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNicId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, datacenterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server ID: %v", serverId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("NIC with ID: %v", nicId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Firewall Rules..."))
	firewallRules, resp, err := c.CloudApiV6Services.FirewallRules().List(datacenterId, serverId, nicId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	firewallRulesItems, ok := firewallRules.GetItemsOk()
	if !ok || firewallRulesItems == nil {
		return fmt.Errorf("could not get items of Firewall Rules")
	}

	if len(*firewallRulesItems) <= 0 {
		return fmt.Errorf("no Firewall Rule found")
	}

	var multiErr error
	for _, firewall := range *firewallRulesItems {
		id := firewall.GetId()
		name := firewall.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Firewall Rule with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.FlagForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, *id, queryParams)

		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

// checkSourceIPAndTargetIPVersions returns true if the source and destination
// IPs are of a different type (IPv4/IPv6) than the specified IP version.
func checkSourceIPAndTargetIPVersions(c *core.CommandConfig) bool {
	ipVersion := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion))

	sIp := net.ParseIP(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagSourceIp)))
	tIp := net.ParseIP(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDestinationIp)))

	isIPv4 := func(ip net.IP) bool {
		return ip != nil && ip.To4() != nil
	}

	if (isIPv4(sIp) || isIPv4(tIp)) && ipVersion == "IPv6" {
		return true
	}

	if (!isIPv4(sIp) || !isIPv4(tIp)) && ipVersion == "IPv4" {
		return true
	}

	return false
}
