package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allFirewallRuleJSONPaths = map[string]string{
		"FirewallRuleId": "id",
		"Name":           "properties.name",
		"Protocol":       "properties.protocol",
		"SourceMac":      "properties.sourceMac",
		"SourceIP":       "properties.sourceIp",
		"DestinationIP":  "properties.destinationIp",
		"PortRangeStart": "properties.portRangeStart",
		"PortRangeEnd":   "properties.portRangeEnd",
		"IcmpCode":       "properties.icmpCode",
		"IcmpType":       "properties.icmpType",
		"Direction":      "properties.type",
		"State":          "metadata.state",
	}

	defaultFirewallRuleCols = []string{"FirewallRuleId", "Name", "Protocol", "PortRangeStart", "PortRangeEnd", "Direction", "State"}
	allFirewallRuleCols     = []string{"FirewallRuleId", "Name", "Protocol", "SourceMac", "SourceIP", "DestinationIP", "PortRangeStart", "PortRangeEnd",
		"IcmpCode", "IcmpType", "Direction", "State"}
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultFirewallRuleCols, printer.ColsMessage(allFirewallRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(firewallRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, "", cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgFirewallRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgServerId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, "", cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Rule", "The name for the Firewall Rule")
	create.AddStringFlag(cloudapiv6.ArgProtocol, "", "", "The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"TCP", "UDP", "ICMP", "ANY"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses")
	create.AddIpFlag(cloudapiv6.ArgSourceIp, "", nil, "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs")
	create.AddIpFlag(cloudapiv6.ArgDestinationIp, cloudapiv6.ArgDestinationIpShort, nil, "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs. WARNING: This short-hand flag `-D` is deprecated.")
	create.AddIntFlag(cloudapiv6.ArgIcmpType, "", 0, "Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types")
	create.AddIntFlag(cloudapiv6.ArgIcmpCode, "", 0, "Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes")
	create.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	create.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	create.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "INGRESS", "The type/direction of Firewall Rule")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule creation [seconds]")
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(cloudapiv6.ArgDepth, "", cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Firewall Rule")
	update.AddStringFlag(cloudapiv6.ArgSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Not setting option allows all source MAC addresses")
	update.AddIpFlag(cloudapiv6.ArgSourceIp, "", nil, "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs")
	update.AddIpFlag(cloudapiv6.ArgDestinationIp, cloudapiv6.ArgDestinationIpShort, nil, "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target/destination IPs. WARNING: This short-hand flag `-D` is deprecated.")
	update.AddIntFlag(cloudapiv6.ArgIcmpType, "", 0, "Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types")
	update.AddIntFlag(cloudapiv6.ArgIcmpCode, "", 0, "Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes")
	update.AddIntFlag(cloudapiv6.ArgPortRangeStart, "", 1, "Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	update.AddIntFlag(cloudapiv6.ArgPortRangeEnd, "", 1, "Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports")
	update.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "", "The type/direction of Firewall Rule")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgFirewallRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgServerId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule update [seconds]")
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(cloudapiv6.ArgDepth, "", cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgNicId, "", "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgFirewallRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.FirewallRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgServerId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Firewall Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Firewalls.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, "", cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return firewallRuleCmd
}

func PreRunFirewallRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.FirewallRulesFilters(), completer.FirewallRulesFiltersUsage())
	}
	return nil
}

func PreRunDcServerNicIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId)
}

func PreRunFirewallDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgAll},
	)
}

func PreRunDcServerNicIdsFRuleProtocol(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgProtocol)
}

func PreRunDcServerNicFRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgFirewallRuleId)
}

func RunFirewallRuleList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	firewallRules, resp, err := c.CloudApiV6Services.FirewallRules().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		listQueryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allFirewallRuleJSONPaths, firewallRules.FirewallRules,
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
		"Firewall Rule with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId))))

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allFirewallRuleJSONPaths, firewallRule.FirewallRule,
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
	properties := getFirewallRulePropertiesSet(c)

	if !properties.HasName() {
		properties.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if !properties.HasType() {
		properties.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))
	}

	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allFirewallRuleJSONPaths, firewallRule.FirewallRule,
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
	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId)),
		getFirewallRulePropertiesSet(c),
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allFirewallRuleJSONPaths, firewallRule.FirewallRule,
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
	datacenterId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgNicId))
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllFirewallRuses(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Firewall Rules successfully deleted"))

		return nil
	}

	if !confirm.Ask("delete firewall rule", viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce))) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Firewall Rule with id: %v...", fruleId))

	resp, err := c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, fruleId, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Firewall Rule successfully deleted"))

	return nil

}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *core.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		properties.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		protocol := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))
		properties.SetProtocol(protocol)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Protocol set: %v", protocol))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceIp)) {
		sourceIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceIp))
		properties.SetSourceIp(sourceIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SourceIp set: %v", sourceIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceMac)) {
		sourceMac := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceMac))
		properties.SetSourceMac(sourceMac)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SourceMac set: %v", sourceMac))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDestinationIp)) {
		targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDestinationIp))
		properties.SetTargetIp(targetIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property TargetIp/DestinationIp set: %v", targetIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpCode)) {
		icmpCode := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpCode))
		properties.SetIcmpCode(icmpCode)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IcmpCode set: %v", icmpCode))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpType)) {
		icmpType := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpType))
		properties.SetIcmpType(icmpType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property IcmpType set: %v", icmpType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart)) {
		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart))
		properties.SetPortRangeStart(portRangeStart)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PortRangeStart set: %v", portRangeStart))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd)) {
		portRangeEnd := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd))
		properties.SetPortRangeEnd(portRangeEnd)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property PortRangeEnd set: %v", portRangeEnd))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)) {
		firewallruleType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection))
		properties.SetType(strings.ToUpper(firewallruleType))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Type/Direction set: %v", firewallruleType))
	}

	return properties
}

func DeleteAllFirewallRuses(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	datacenterId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.Resource, cloudapiv6.ArgNicId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Datacenter ID: %v", datacenterId))
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Firewall Rules to be deleted:"))

	for _, firewall := range *firewallRulesItems {
		delIdAndName := ""

		if id, ok := firewall.GetIdOk(); ok && id != nil {
			delIdAndName += "Firewall Rule Id: " + *id
		}

		if properties, ok := firewall.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " Firewall Rule Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Firewall Rules", viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce))) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Firewall Rules..."))

	var multiErr error
	for _, firewall := range *firewallRulesItems {
		id, ok := firewall.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Firewall Rule with id: %v...", *id))

		resp, err = c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, *id, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}
	return nil
}
