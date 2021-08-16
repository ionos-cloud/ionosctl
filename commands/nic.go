package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

func nic() *core.Command {
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
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = nicCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", "The unique Server Id")
	_ = viper.BindPFlag(core.GetGlobalFlagName(nicCmd.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = nicCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSliceP(config.ArgCols, "", defaultNicCols, utils.ColsMessage(allNicCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nicCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nicCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, nicCmd, core.CommandBuilder{
		Namespace:  "nic",
		Resource:   "nic",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List NICs",
		LongDesc:   "Use this command to get a list of NICs on your account.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listNicExample,
		PreCmdRun:  PreRunGlobalDcServerIds,
		CmdRun:     RunNicList,
		InitClient: true,
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
		PreCmdRun:  PreRunGlobalDcServerIdsNicId,
		CmdRun:     RunNicGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId)), viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunGlobalDcServerIds,
		CmdRun:     RunNicCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Internet Access", "The name of the NIC")
	create.AddStringSliceFlag(config.ArgIps, "", []string{""}, "IPs assigned to the NIC. This can be a collection")
	create.AddBoolFlag(config.ArgDhcp, "", config.DefaultDhcp, "Set to false if you wish to disable DHCP on the NIC")
	create.AddIntFlag(config.ArgLanId, "", config.DefaultNicLanId, "The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunGlobalDcServerIdsNicId,
		CmdRun:     RunNicUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId)), viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name of the NIC")
	update.AddIntFlag(config.ArgLanId, "", config.DefaultNicLanId, "The LAN ID the NIC sits on")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgDhcp, "", config.DefaultDhcp, "Boolean value that indicates if the NIC is using DHCP (true) or not (false)")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC update [seconds]")
	update.AddStringSliceFlag(config.ArgIps, "", []string{""}, "IPs assigned to the NIC")

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
		PreCmdRun:  PreRunGlobalDcServerIdsNicId,
		CmdRun:     RunNicDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgDataCenterId)), viper.GetString(core.GetGlobalFlagName(nicCmd.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC deletion [seconds]")

	return nicCmd
}

func PreRunGlobalDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgServerId)
}

func PreRunGlobalDcServerIdsNicId(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgServerId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgNicId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunNicList(c *core.CommandConfig) error {
	nics, resp, err := c.Nics().List(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	ss := getNics(nics)
	return c.Printer.Print(printer.Result{
		OutputJSON: nics,
		KeyValue:   getNicsKVMaps(ss),
		Columns:    getNicsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunNicGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Nic with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)))
	nic, resp, err := c.Nics().Get(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, []v5.Nic{*nic}))
}

func RunNicCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId))
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	ips := viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgIps))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, config.ArgDhcp))
	lanId := viper.GetInt32(core.GetFlagName(c.NS, config.ArgLanId))

	c.Printer.Verbose("Creating Nic in DataCenterId: %v with ServerId: %v...", dcId, serverId)
	c.Printer.Verbose("Properties set for creating the Nic: Name: %v, Ips: %v, Dhcp: %v, Lan: %v",
		name, ips, dhcp, lanId)
	nic, resp, err := c.Nics().Create(dcId, serverId, name, ips, dhcp, lanId)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}

	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, []v5.Nic{*nic}))
}

func RunNicUpdate(c *core.CommandConfig) error {
	input := v5.NicProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
		input.NicProperties.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, config.ArgDhcp))
		input.NicProperties.SetDhcp(dhcp)
		c.Printer.Verbose("Property Dhcp set: %v", dhcp)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgLanId)) {
		lan := viper.GetInt32(core.GetFlagName(c.NS, config.ArgLanId))
		input.NicProperties.SetLan(lan)
		c.Printer.Verbose("Property Lan set: %v", lan)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIps)) {
		ips := viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgIps))
		input.NicProperties.SetIps(ips)
		c.Printer.Verbose("Property Ips set: %v", ips)
	}
	nicUpd, resp, err := c.Nics().Update(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, []v5.Nic{*nicUpd}))
}

func RunNicDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nic"); err != nil {
		return err
	}
	c.Printer.Verbose("nic with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)))
	resp, err := c.Nics().Delete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, nil))
}

// LoadBalancer Nic Commands

func loadBalancerNic() *core.Command {
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
	attachNic.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, utils.ColsMessage(allNicCols))
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgServerId, "", "", "The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion")
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachNic.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(attachNic.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = attachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(attachNic.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(attachNic.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
	listNics.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, utils.ColsMessage(allNicCols))
	_ = listNics.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	listNics.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = listNics.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listNics.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = listNics.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(listNics.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	getNicCmd.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, utils.ColsMessage(allNicCols))
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(getNicCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getNicCmd.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = getNicCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(os.Stderr,
			viper.GetString(core.GetFlagName(getNicCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(getNicCmd.NS, config.ArgLoadBalancerId)),
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
	detachNic.AddStringSliceFlag(config.ArgCols, "", defaultNicCols, utils.ColsMessage(allNicCols))
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNicCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(config.ArgNicId, config.ArgIdShort, "", config.RequiredFlagNicId)
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(os.Stderr, viper.GetString(core.GetFlagName(detachNic.NS, config.ArgDataCenterId)), viper.GetString(core.GetFlagName(detachNic.NS, config.ArgLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	_ = detachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, viper.GetString(core.GetFlagName(detachNic.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NIC detachment to be executed")
	detachNic.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NIC detachment [seconds]")

	return loadbalancerNicCmd
}

func PreRunDcNicLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNicId, config.ArgLoadBalancerId)
}

func RunLoadBalancerNicAttach(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	lbId := viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, config.ArgNicId))
	c.Printer.Verbose("Attaching Nic with id: %v from LoadBalancer with id: %v")
	attachedNic, resp, err := c.Loadbalancers().AttachNic(dcId, lbId, nicId)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, getNic(attachedNic)))
}

func RunLoadBalancerNicList(c *core.CommandConfig) error {
	attachedNics, resp, err := c.Loadbalancers().ListNics(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getAttachedNics(attachedNics)))
}

func RunLoadBalancerNicGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Nic with id: %v from LoadBalancer with id: %v...")
	n, resp, err := c.Loadbalancers().GetNic(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(nil, c, getNic(n)))
}

func RunLoadBalancerNicDetach(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	lbId := viper.GetString(core.GetFlagName(c.NS, config.ArgLoadBalancerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, config.ArgNicId))
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach nic from loadbalancer"); err != nil {
		return err
	}
	c.Printer.Verbose("Detaching Nic with id: %v from LoadBalancer with id: %v...")
	resp, err := c.Loadbalancers().DetachNic(dcId, lbId, nicId)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNicPrint(resp, c, nil))
}

// Output Printing

var (
	defaultNicCols = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "State"}
	allNicCols     = []string{"NicId", "Name", "Dhcp", "LanId", "Ips", "State", "FirewallActive", "Mac"}
)

type NicPrint struct {
	NicId          string   `json:"NicId,omitempty"`
	Name           string   `json:"Name,omitempty"`
	Dhcp           bool     `json:"Dhcp,omitempty"`
	LanId          int32    `json:"LanId,omitempty"`
	Ips            []string `json:"Ips,omitempty"`
	FirewallActive bool     `json:"FirewallActive,omitempty"`
	Mac            string   `json:"Mac,omitempty"`
	State          string   `json:"State,omitempty"`
}

func getNicPrint(resp *v5.Response, c *core.CommandConfig, nics []v5.Nic) printer.Result {
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
		"Mac":            "Mac",
		"State":          "State",
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

func getNics(nics v5.Nics) []v5.Nic {
	ns := make([]v5.Nic, 0)
	for _, s := range *nics.Items {
		ns = append(ns, v5.Nic{Nic: s})
	}
	return ns
}

func getNic(n *v5.Nic) []v5.Nic {
	nics := make([]v5.Nic, 0)
	if n != nil {
		nics = append(nics, v5.Nic{Nic: n.Nic})
	}
	return nics
}

func getAttachedNics(nics v5.BalancedNics) []v5.Nic {
	ns := make([]v5.Nic, 0)
	for _, s := range *nics.BalancedNics.Items {
		ns = append(ns, v5.Nic{Nic: s})
	}
	return ns
}

func getNicsKVMaps(ns []v5.Nic) []map[string]interface{} {
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

func getNicsIds(outErr io.Writer, datacenterId, serverId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	nicSvc := v5.NewNicService(clientSvc.Get(), context.TODO())
	nics, _, err := nicSvc.List(datacenterId, serverId)
	clierror.CheckError(err, outErr)
	nicsIds := make([]string, 0)
	if items, ok := nics.Nics.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				nicsIds = append(nicsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return nicsIds
}

func getAttachedNicsIds(outErr io.Writer, datacenterId, loadbalancerId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	nicSvc := v5.NewLoadbalancerService(clientSvc.Get(), context.TODO())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId)
	clierror.CheckError(err, outErr)
	attachedNicsIds := make([]string, 0)
	if items, ok := nics.BalancedNics.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedNicsIds = append(attachedNicsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedNicsIds
}
