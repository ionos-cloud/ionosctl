package commands

import (
	"context"
	"errors"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NicCmd() *core.Command {
	ctx := context.TODO()
	nicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Network Interfaces Operations",
			Long:             "The sub-commands of `ionosctl nic` allow you to create, list, get, update, delete NICs. To attach a NIC to a Load Balancer, use the Load Balancer command `ionosctl loadbalancer nic attach`.",
			TraverseChildren: true,
		},
	}
	globalFlags := nicCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultNicCols, printer.ColsMessage(allNicCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nicCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nicCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace:  "nic",
		Resource:   "nic",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List NICs",
		LongDesc:   "Use this command to get a list of NICs on your account.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listNicExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunNicList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace:  "nic",
		Resource:   "nic",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a NIC",
		LongDesc:   "Use this command to get information about a specified NIC from specified Data Center and Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* NIC Id",
		Example:    getNicExample,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunNicGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a NIC",
		LongDesc: `Use this command to create/add a new NIC to the target Server. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run a command:

* Data Center Id
* Server Id`,
		Example:    createNicExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunNicCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "Internet Access", "The name of the NIC")
	create.AddStringSliceFlag(cloudapi_v6.ArgIps, "", []string{""}, "IPs assigned to the NIC. This can be a collection")
	create.AddBoolFlag(cloudapi_v6.ArgDhcp, "", cloudapi_v6.DefaultDhcp, "Set to false if you wish to disable DHCP on the NIC")
	create.AddBoolFlag(cloudapi_v6.ArgFirewallActive, "", cloudapi_v6.DefaultFirewallActive, "Activate or deactivate the Firewall")
	create.AddStringFlag(cloudapi_v6.ArgFirewallType, "", "INGRESS", "The type of Firewall Rules that will be allowed on the NIC")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFirewallType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapi_v6.ArgLanId, "", cloudapi_v6.DefaultNicLanId, "The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a NIC",
		LongDesc: `Use this command to update the configuration of a specified NIC. Some restrictions are in place: The primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer. You can also add additional reserved, public IPs to the NIC.

The user can specify and assign private IPs manually. Valid IP addresses for private networks are 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16.

The value for firewallActive can be toggled between true and false to enable or disable the firewall. When the firewall is enabled, incoming traffic is filtered by all the firewall rules configured on the NIC. When the firewall is disabled, then all incoming traffic is routed directly to the NIC and any configured firewall rules are ignored.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* NIC Id`,
		Example:    updateNicExample,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunNicUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "", "The name of the NIC")
	update.AddIntFlag(cloudapi_v6.ArgLanId, "", cloudapi_v6.DefaultNicLanId, "The LAN ID the NIC sits on")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapi_v6.ArgFirewallActive, "", cloudapi_v6.DefaultFirewallActive, "Activate or deactivate the Firewall")
	update.AddStringFlag(cloudapi_v6.ArgFirewallType, "", "INGRESS", "The type of Firewall Rules that will be allowed on the NIC")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFirewallType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapi_v6.ArgDhcp, "", cloudapi_v6.DefaultDhcp, "Boolean value that indicates if the NIC is using DHCP (true) or not (false)")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC update [seconds]")
	update.AddStringSliceFlag(cloudapi_v6.ArgIps, "", []string{""}, "IPs assigned to the NIC")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace: "nic",
		Resource:  "nic",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a NIC",
		LongDesc: `This command deletes a specified NIC.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id`,
		Example:    deleteNicExample,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunNicDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC deletion [seconds]")

	return nicCmd
}

func RunNicList(c *core.CommandConfig) error {
	nics, _, err := c.CloudApiV6Services.Nics().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getNics(nics)))
}

func RunNicGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Nic with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)))
	n, _, err := c.CloudApiV6Services.Nics().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, []resources.Nic{*n}))
}

func RunNicCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
	ips := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapi_v6.ArgIps))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDhcp))
	lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapi_v6.ArgLanId))
	firewallActive := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallActive))
	firewallType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallType))

	c.Printer.Verbose("Creating Nic in DataCenterId: %v with ServerId: %v...", dcId, serverId)
	c.Printer.Verbose("Properties set for creating the Nic: Name: %v, Ips: %v, Dhcp: %v, Lan: %v FirewallActive: %v, FirewallType: %v",
		name, ips, dhcp, lanId, firewallActive, firewallType)

	inputProper := resources.NicProperties{}
	inputProper.SetName(name)
	inputProper.SetIps(ips)
	inputProper.SetDhcp(dhcp)
	inputProper.SetLan(lanId)
	inputProper.SetFirewallActive(firewallActive)
	inputProper.SetFirewallType(firewallType)
	input := resources.Nic{
		Nic: ionoscloud.Nic{
			Properties: &inputProper.NicProperties,
		},
	}
	n, resp, err := c.CloudApiV6Services.Nics().Create(dcId, serverId, input)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, []resources.Nic{*n}))
}

func RunNicUpdate(c *core.CommandConfig) error {
	input := resources.NicProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
		input.NicProperties.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDhcp))
		input.NicProperties.SetDhcp(dhcp)
		c.Printer.Verbose("Property Dhcp set: %v", dhcp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgLanId)) {
		lan := viper.GetInt32(core.GetFlagName(c.NS, cloudapi_v6.ArgLanId))
		input.NicProperties.SetLan(lan)
		c.Printer.Verbose("Property Lan set: %v", lan)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgIps)) {
		ips := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapi_v6.ArgIps))
		input.NicProperties.SetIps(ips)
		c.Printer.Verbose("Property Ips set: %v", ips)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallActive)) {
		firewallActive := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallActive))
		input.NicProperties.SetFirewallActive(firewallActive)
		c.Printer.Verbose("Property FirewallActive set: %v", firewallActive)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallType)) {
		firewallType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFirewallType))
		input.NicProperties.SetFirewallType(firewallType)
		c.Printer.Verbose("Property FirewallType set: %v", firewallType)
	}
	nicUpd, resp, err := c.CloudApiV6Services.Nics().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, []resources.Nic{*nicUpd}))
}

func RunNicDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nic"); err != nil {
		return err
	}
	c.Printer.Verbose("Nic with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)))
	resp, err := c.CloudApiV6Services.Nics().Delete(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, nil))
}

// LoadBalancer Nic Commands

func LoadBalancerNicCmd() *core.Command {
	ctx := context.TODO()
	loadbalancerNicCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Aliases:          []string{"n"},
			Short:            "Load Balancer Nic Operations",
			Long:             "The sub-commands of `ionosctl loadbalancer nic` allow you to manage NICs on Load Balancers on your account.",
			TraverseChildren: true,
		},
	}

	/*
		Attach Nic Command
	*/
	attachNic := core.NewCommand(ctx, loadbalancerNicCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "nic",
		Verb:      "attach",
		Aliases:   []string{"a"},
		ShortDesc: "Attach a NIC to a Load Balancer",
		LongDesc: `Use this command to associate a NIC to a Load Balancer, enabling the NIC to participate in load-balancing.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`,
		Example:    attachNicLoadbalancerExample,
		PreCmdRun:  PreRunDcNicLoadBalancerIds,
		CmdRun:     RunLoadBalancerNicAttach,
		InitClient: true,
	})
	attachNic.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, printer.ColsMessage(allNicCols))
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = attachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(cloudapi_v6.ArgServerId, "", "", "The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion")
	_ = attachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachNic.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(cloudapi_v6.ArgLoadBalancerId, "", "", cloudapi_v6.LoadBalancerId, core.RequiredFlagOption())
	_ = attachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(attachNic.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = attachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(attachNic.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(attachNic.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC attachment to be executed")
	attachNic.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC attachment [seconds]")

	/*
		List Nics Command
	*/
	listNics := core.NewCommand(ctx, loadbalancerNicCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "nic",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List attached NICs from a Load Balancer",
		LongDesc:   "Use this command to get a list of attached NICs to a Load Balancer from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		Example:    listNicsLoadbalancerExample,
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerNicList,
		InitClient: true,
	})
	listNics.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, printer.ColsMessage(allNicCols))
	_ = listNics.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	listNics.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = listNics.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listNics.AddStringFlag(cloudapi_v6.ArgLoadBalancerId, "", "", cloudapi_v6.LoadBalancerId, core.RequiredFlagOption())
	_ = listNics.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(listNics.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Nic Command
	*/
	getNicCmd := core.NewCommand(ctx, loadbalancerNicCmd, core.CommandBuilder{
		Namespace:  "loadbalancer",
		Resource:   "nic",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an attached NIC to a Load Balancer",
		LongDesc:   "Use this command to retrieve the attributes of a given load balanced NIC.\n\nRequired values to run the command:\n\n* Data Center Id\n* Load Balancer Id\n* NIC Id",
		Example:    getNicLoadbalancerExample,
		PreCmdRun:  PreRunDcNicLoadBalancerIds,
		CmdRun:     RunLoadBalancerNicGet,
		InitClient: true,
	})
	getNicCmd.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, printer.ColsMessage(allNicCols))
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(cloudapi_v6.ArgLoadBalancerId, "", "", cloudapi_v6.LoadBalancerId, core.RequiredFlagOption())
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(getNicCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedNicsIds(os.Stderr, viper.GetString(core.GetFlagName(getNicCmd.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(getNicCmd.NS, cloudapi_v6.ArgLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Nic Command
	*/
	detachNic := core.NewCommand(ctx, loadbalancerNicCmd, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "nic",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a NIC from a Load Balancer",
		LongDesc: `Use this command to remove the association of a NIC with a Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`,
		Example:    detachNicLoadbalancerExample,
		PreCmdRun:  PreRunDcNicLoadBalancerIds,
		CmdRun:     RunLoadBalancerNicDetach,
		InitClient: true,
	})
	detachNic.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, printer.ColsMessage(allNicCols))
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = detachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(cloudapi_v6.ArgNicId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = detachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedNicsIds(os.Stderr, viper.GetString(core.GetFlagName(detachNic.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachNic.NS, cloudapi_v6.ArgLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(cloudapi_v6.ArgLoadBalancerId, "", "", cloudapi_v6.LoadBalancerId, core.RequiredFlagOption())
	_ = detachNic.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(detachNic.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC detachment to be executed")
	detachNic.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC detachment [seconds]")

	return loadbalancerNicCmd
}

func PreRunDcNicLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgNicId, cloudapi_v6.ArgLoadBalancerId)
}

func RunLoadBalancerNicAttach(c *core.CommandConfig) error {
	attachedNic, resp, err := c.CloudApiV6Services.Loadbalancers().AttachNic(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, getNic(attachedNic)))
}

func RunLoadBalancerNicList(c *core.CommandConfig) error {
	attachedNics, _, err := c.CloudApiV6Services.Loadbalancers().ListNics(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getAttachedNics(attachedNics)))
}

func RunLoadBalancerNicGet(c *core.CommandConfig) error {
	n, _, err := c.CloudApiV6Services.Loadbalancers().GetNic(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getNic(n)))
}

func RunLoadBalancerNicDetach(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach nic from loadbalancer"); err != nil {
		return err
	}
	resp, err := c.CloudApiV6Services.Loadbalancers().DetachNic(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, nil))
}

// Output Printing

var (
	defaultNicCols = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "State"}
	allNicCols     = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "State", "FirewallActive", "FirewallType", "DeviceNumber", "PciSlot", "Mac"}
)

type NicPrint struct {
	NicId          string   `json:"NicId,omitempty"`
	Name           string   `json:"Name,omitempty"`
	Dhcp           bool     `json:"Dhcp,omitempty"`
	LanId          int32    `json:"LanId,omitempty"`
	Ips            []string `json:"Ips,omitempty"`
	FirewallActive bool     `json:"FirewallActive,omitempty"`
	FirewallType   string   `json:"FirewallType,omitempty"`
	Mac            string   `json:"Mac,omitempty"`
	State          string   `json:"State,omitempty"`
	DeviceNumber   int32    `json:"DeviceNumber,omitempty"`
	PciSlot        int32    `json:"PciSlot,omitempty"`
}

func getNicPrint(resp *resources.Response, c *core.CommandConfig, nics []resources.Nic) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if nics != nil {
			r.OutputJSON = nics
			r.KeyValue = getNicsKVMaps(nics)
			if c.Resource != c.Namespace {
				r.Columns = getNicsCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getNicsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
		}
	}
	return r
}

func getNicsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultNicCols
	}

	columnsMap := map[string]string{
		"NicId":          "NicId",
		"Name":           "Name",
		"Dhcp":           "Dhcp",
		"LanId":          "LanId",
		"Ips":            "Ips",
		"FirewallActive": "FirewallActive",
		"FirewallType":   "FirewallType",
		"Mac":            "Mac",
		"State":          "State",
		"DeviceNumber":   "DeviceNumber",
		"PciSlot":        "PciSlot",
	}
	var nicCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			nicCols = append(nicCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return nicCols
}

func getNics(nics resources.Nics) []resources.Nic {
	ns := make([]resources.Nic, 0)
	for _, s := range *nics.Items {
		ns = append(ns, resources.Nic{Nic: s})
	}
	return ns
}

func getNic(n *resources.Nic) []resources.Nic {
	nics := make([]resources.Nic, 0)
	if n != nil {
		nics = append(nics, resources.Nic{Nic: n.Nic})
	}
	return nics
}

func getAttachedNics(nics resources.BalancedNics) []resources.Nic {
	ns := make([]resources.Nic, 0)
	for _, s := range *nics.BalancedNics.Items {
		ns = append(ns, resources.Nic{Nic: s})
	}
	return ns
}

func getNicsKVMaps(ns []resources.Nic) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ns))
	for _, n := range ns {
		var nicprint NicPrint
		if id, ok := n.GetIdOk(); ok && id != nil {
			nicprint.NicId = *id
		}
		if properties, ok := n.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				nicprint.Name = *name
			}
			if dhcp, ok := properties.GetDhcpOk(); ok && dhcp != nil {
				nicprint.Dhcp = *dhcp
			}
			if lanId, ok := properties.GetLanOk(); ok && lanId != nil {
				nicprint.LanId = *lanId
			}
			if factive, ok := properties.GetFirewallActiveOk(); ok && factive != nil {
				nicprint.FirewallActive = *factive
			}
			if ftype, ok := properties.GetFirewallTypeOk(); ok && ftype != nil {
				nicprint.FirewallType = *ftype
			}
			if no, ok := properties.GetDeviceNumberOk(); ok && no != nil {
				nicprint.DeviceNumber = *no
			}
			if slot, ok := properties.GetPciSlotOk(); ok && slot != nil {
				nicprint.PciSlot = *slot
			}
			if mac, ok := properties.GetMacOk(); ok && mac != nil {
				nicprint.Mac = *mac
			}
			if ips, ok := properties.GetIpsOk(); ok && ips != nil {
				nicprint.Ips = *ips
			}
		}
		if metadata, ok := n.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				nicprint.State = *state
			}
		}
		o := structs.Map(nicprint)
		out = append(out, o)
	}
	return out
}
