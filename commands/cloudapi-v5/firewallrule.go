package cloudapi_v5

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FirewallRuleCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultFirewallRuleCols, printer.ColsMessage(allFirewallRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(firewallRuleCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		LongDesc:   "Use this command to get a list of Firewall Rules from a specified NIC from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Nic Id",
		Example:    listFirewallRuleExample,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunFirewallRuleList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgNicId, "", "", cloudapiv5.NicId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgNicId, "", "", cloudapiv5.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgFirewallRuleId, cloudapiv5.ArgIdShort, "", cloudapiv5.FirewallRuleId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgServerId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})

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
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "Unnamed Rule", "The name for the Firewall Rule")
	create.AddStringFlag(cloudapiv5.ArgProtocol, "", "", "The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY", core.RequiredFlagOption())
	create.AddStringFlag(cloudapiv5.ArgSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses.")
	create.AddStringFlag(cloudapiv5.ArgSourceIp, "", "", "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs.")
	create.AddStringFlag(cloudapiv5.ArgTargetIp, "", "", "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target IPs.")
	create.AddIntFlag(cloudapiv5.ArgIcmpType, "", 0, "Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types.")
	create.AddIntFlag(cloudapiv5.ArgIcmpCode, "", 0, "Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes.")
	create.AddIntFlag(cloudapiv5.ArgPortRangeStart, "", 1, "Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports.")
	create.AddIntFlag(cloudapiv5.ArgPortRangeStop, "", 1, "Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports.")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Firewall Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule creation [seconds]")
	create.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv5.ArgNicId, "", "", cloudapiv5.NicId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(create.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(create.NS, cloudapiv5.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

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
	update.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgNicId, "", "", cloudapiv5.NicId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "The name for the Firewall Rule")
	update.AddStringFlag(cloudapiv5.ArgSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Not setting option allows all source MAC addresses.")
	update.AddStringFlag(cloudapiv5.ArgSourceIp, "", "", "Only traffic originating from the respective IPv4 address is allowed. Not setting option allows all source IPs.")
	update.AddStringFlag(cloudapiv5.ArgTargetIp, "", "", "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Not setting option allows all target IPs.")
	update.AddIntFlag(cloudapiv5.ArgIcmpType, "", 0, "Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Not setting option allows all types.")
	update.AddIntFlag(cloudapiv5.ArgIcmpCode, "", 0, "Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Not setting option allows all codes.")
	update.AddIntFlag(cloudapiv5.ArgPortRangeStart, "", 1, "Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports.")
	update.AddIntFlag(cloudapiv5.ArgPortRangeStop, "", 1, "Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Not setting portRangeStart and portRangeEnd allows all ports.")
	update.AddStringFlag(cloudapiv5.ArgFirewallRuleId, cloudapiv5.ArgIdShort, "", cloudapiv5.FirewallRuleId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgServerId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapiv5.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Firewall Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule update [seconds]")

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
		PreCmdRun:  PreRunDcServerNicDelete,
		CmdRun:     RunFirewallRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgDataCenterId, "", "", cloudapiv5.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgServerId, "", "", cloudapiv5.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgNicId, "", "", cloudapiv5.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgFirewallRuleId, cloudapiv5.ArgIdShort, "", cloudapiv5.FirewallRuleId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FirewallRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgServerId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for Firewall Rule deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Firewalls.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Firewall Rule deletion [seconds]")

	return firewallRuleCmd
}

func PreRunDcServerNicIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgNicId)
}

func PreRunDcServerNicDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgNicId, cloudapiv5.ArgFirewallRuleId},
		[]string{cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgNicId, cloudapiv5.ArgAll},
	)
}

func PreRunDcServerNicIdsFRuleProtocol(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgNicId, cloudapiv5.ArgProtocol)
}

func PreRunDcServerNicFRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgDataCenterId, cloudapiv5.ArgServerId, cloudapiv5.ArgNicId, cloudapiv5.ArgFirewallRuleId)
}

func RunFirewallRuleList(c *core.CommandConfig) error {
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgNicId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	c.Printer.Verbose("Getting Firewall Rules from NIC with ID: %v; Server ID: %v; Datacenter ID: %v... ", nicId, serverId, datacenterId)
	firewallRules, resp, err := c.CloudApiV5Services.FirewallRules().List(datacenterId, serverId, nicId)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(nil, c, getFirewallRules(firewallRules)))
}

func RunFirewallRuleGet(c *core.CommandConfig) error {
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgFirewallRuleId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgNicId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	c.Printer.Verbose("Getting Firewall Rule with ID: %v from NIC with ID: %v; Server ID: %v; Datacenter ID: %v... ", fruleId, nicId, serverId, datacenterId)
	firewallRule, resp, err := c.CloudApiV5Services.FirewallRules().Get(datacenterId, serverId, nicId, fruleId)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(nil, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleCreate(c *core.CommandConfig) error {
	properties := getFirewallRulePropertiesSet(c)
	if !properties.HasName() {
		properties.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName)))
	}
	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgNicId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	c.Printer.Verbose("Creating Firewall Rule attached to specified NIC with ID: %v; Server ID: %v; Datacenter ID: %v... ", nicId, serverId, datacenterId)
	firewallRule, resp, err := c.CloudApiV5Services.FirewallRules().Create(datacenterId, serverId, nicId, input)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(resp, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleUpdate(c *core.CommandConfig) error {
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgFirewallRuleId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgNicId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	c.Printer.Verbose("Updating Firewall Rule with ID: %v from NIC with ID: %v; Server ID: %v; Datacenter ID: %v... ", fruleId, nicId, serverId, datacenterId)
	firewallRule, resp, err := c.CloudApiV5Services.FirewallRules().Update(datacenterId, serverId, nicId, fruleId, getFirewallRulePropertiesSet(c))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(resp, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	datacenterId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgServerId))
	nicId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgNicId))
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgFirewallRuleId))
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		err := DeleteAllFirewallRules(c)
		if err != nil {
			return err
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete firewall rule"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting Firewall Rule with ID: %v from NIC with ID: %v; Server ID: %v; Datacenter ID: %v... ", fruleId, nicId, serverId, datacenterId)
		resp, err := c.CloudApiV5Services.FirewallRules().Delete(datacenterId, serverId, nicId, fruleId)
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

	return c.Printer.Print(getFirewallRulePrint(resp, c, nil))
}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *core.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
		properties.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgProtocol)) {
		protocol := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgProtocol))
		properties.SetProtocol(protocol)
		c.Printer.Verbose("Property Protocol set: %v", protocol)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgSourceIp)) {
		sourceIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSourceIp))
		properties.SetSourceIp(sourceIp)
		c.Printer.Verbose("Property SourceIp set: %v", sourceIp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgSourceMac)) {
		sourceMac := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSourceMac))
		properties.SetSourceMac(sourceMac)
		c.Printer.Verbose("Property SourceMac set: %v", sourceMac)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgTargetIp)) {
		targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgTargetIp))
		properties.SetTargetIp(targetIp)
		c.Printer.Verbose("Property TargetIp set: %v", targetIp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgIcmpCode)) {
		icmpCode := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgIcmpCode))
		properties.SetIcmpCode(icmpCode)
		c.Printer.Verbose("Property IcmpCode set: %v", icmpCode)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgIcmpType)) {
		icmpType := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgIcmpType))
		properties.SetIcmpType(icmpType)
		c.Printer.Verbose("Property IcmpType set: %v", icmpType)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgPortRangeStart)) {
		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgPortRangeStart))
		properties.SetPortRangeStart(portRangeStart)
		c.Printer.Verbose("Property PortRangeStart set: %v", portRangeStart)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgPortRangeStop)) {
		portRangeStop := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv5.ArgPortRangeStop))
		properties.SetPortRangeEnd(portRangeStop)
		c.Printer.Verbose("Property PortRangeEnd set: %v", portRangeStop)
	}
	return properties
}

func DeleteAllFirewallRules(c *core.CommandConfig) error {
	datacenterId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgServerId))
	nicId := viper.GetString(core.GetGlobalFlagName(c.Resource, cloudapiv5.ArgNicId))
	_ = c.Printer.Print("Firewallrules to be deleted:")
	firewallrules, resp, err := c.CloudApiV5Services.FirewallRules().List(datacenterId, serverId, nicId)
	if err != nil {
		return err
	}
	if firewallrulestems, ok := firewallrules.GetItemsOk(); ok && firewallrulestems != nil {
		for _, firewall := range *firewallrulestems {
			if id, ok := firewall.GetIdOk(); ok && id != nil {
				_ = c.Printer.Print("Firewallrule Id: " + *id)
			}
			if properties, ok := firewall.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print(" Firewallrule Name: " + *name)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Firewallrules"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the Firewallrules...")
		for _, firewall := range *firewallrulestems {
			if id, ok := firewall.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Deleting Firewall Rule with id: %v...", *id)
				resp, err = c.CloudApiV5Services.FirewallRules().Delete(datacenterId, serverId, nicId, *id)
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
		}
	}
	return nil
}

// Output Printing

var (
	defaultFirewallRuleCols = []string{"FirewallRuleId", "Name", "Protocol", "PortRangeStart", "PortRangeEnd", "State"}
	allFirewallRuleCols     = []string{"FirewallRuleId", "Name", "Protocol", "SourceMac", "SourceIP", "TargetIP", "PortRangeStart", "PortRangeEnd", "State"}
)

type FirewallRulePrint struct {
	FirewallRuleId string `json:"FirewallRuleId,omitempty"`
	Name           string `json:"Name,omitempty"`
	Protocol       string `json:"Protocol,omitempty"`
	SourceMac      string `json:"SourceMac,omitempty"`
	SourceIP       string `json:"SourceIP,omitempty"`
	TargetIP       string `json:"TargetIP,omitempty"`
	PortRangeStart int32  `json:"PortRangeStart,omitempty"`
	PortRangeEnd   int32  `json:"PortRangeEnd,omitempty"`
	State          string `json:"State,omitempty"`
}

func getFirewallRulePrint(resp *resources.Response, c *core.CommandConfig, rule []resources.FirewallRule) printer.Result {
	var r printer.Result
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if rule != nil {
			r.OutputJSON = rule
			r.KeyValue = getFirewallRulesKVMaps(rule)
			r.Columns = getFirewallRulesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getFirewallRulesCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) && len(viper.GetStringSlice(flagName)) > 0 {
		var firewallRuleCols []string
		columnsMap := map[string]string{
			"FirewallRuleId": "FirewallRuleId",
			"Name":           "Name",
			"Protocol":       "Protocol",
			"SourceMac":      "SourceMac",
			"SourceIP":       "SourceIP",
			"TargetIP":       "TargetIP",
			"PortRangeStart": "PortRangeStart",
			"PortRangeEnd":   "PortRangeEnd",
			"State":          "State",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				firewallRuleCols = append(firewallRuleCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return firewallRuleCols
	} else {
		return defaultFirewallRuleCols
	}
}

func getFirewallRules(firewallRules resources.FirewallRules) []resources.FirewallRule {
	ls := make([]resources.FirewallRule, 0)
	if items, ok := firewallRules.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ls = append(ls, resources.FirewallRule{FirewallRule: s})
		}
	}
	return ls
}

func getFirewallRule(firewallRule *resources.FirewallRule) []resources.FirewallRule {
	ss := make([]resources.FirewallRule, 0)
	if firewallRule != nil {
		ss = append(ss, resources.FirewallRule{FirewallRule: firewallRule.FirewallRule})
	}
	return ss
}

func getFirewallRulesKVMaps(ls []resources.FirewallRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	if len(ls) > 0 {
		for _, l := range ls {
			o := getFirewallRuleKVMap(l)
			out = append(out, o)
		}
	}
	return out
}

func getFirewallRuleKVMap(l resources.FirewallRule) map[string]interface{} {
	var firewallRulePrint FirewallRulePrint
	if id, ok := l.GetIdOk(); ok && id != nil {
		firewallRulePrint.FirewallRuleId = *id
	}
	if properties, ok := l.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			firewallRulePrint.Name = *name
		}
		if protocol, ok := properties.GetProtocolOk(); ok && protocol != nil {
			firewallRulePrint.Protocol = *protocol
		}
		if portRangeStart, ok := properties.GetPortRangeStartOk(); ok && portRangeStart != nil {
			firewallRulePrint.PortRangeStart = *portRangeStart
		}
		if portRangeEnd, ok := properties.GetPortRangeEndOk(); ok && portRangeEnd != nil {
			firewallRulePrint.PortRangeEnd = *portRangeEnd
		}
		if sourceMac, ok := properties.GetSourceMacOk(); ok && sourceMac != nil {
			firewallRulePrint.SourceMac = *sourceMac
		}
		if sourceIp, ok := properties.GetSourceIpOk(); ok && sourceIp != nil {
			firewallRulePrint.SourceIP = *sourceIp
		}
		if targetIp, ok := properties.GetTargetIpOk(); ok && targetIp != nil {
			firewallRulePrint.TargetIP = *targetIp
		}
	}
	if metadata, ok := l.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			firewallRulePrint.State = *state
		}
	}
	return structs.Map(firewallRulePrint)
}
