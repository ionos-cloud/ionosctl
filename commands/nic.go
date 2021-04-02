package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func nic() *builder.Command {
	nicCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "nic",
			Short:            "Network Interfaces Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl nic` + "`" + ` allow you to create, list, get, update, delete NICs or attach, detach a NIC from a Load Balancer.`,
			TraverseChildren: true,
		},
	}
	globalFlags := nicCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	viper.BindPFlag(builder.GetGlobalFlagName(nicCmd.Command.Use, config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	nicCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", "The unique Server Id")
	viper.BindPFlag(builder.GetGlobalFlagName(nicCmd.Command.Use, config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	nicCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultNicCols, "Columns to be printed in the standard output")
	viper.BindPFlag(builder.GetGlobalFlagName(nicCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcServerIdsValidate, RunNicList, "list", "List NICs",
		"Use this command to get a list of NICs on your account.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		listNicExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcServerIdsNicIdValidate, RunNicGet, "get", "Get a NIC",
		"Use this command to get information about a specified NIC from specified Data Center and Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* NIC Id",
		getNicExample, true)
	get.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	get.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcServerIdsValidate, RunNicCreate, "create", "Create a NIC",
		`Use this command to create a new NIC on your account. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run a command:

* Data Center Id
* Server Id`, createNicExample, true)
	create.AddStringFlag(config.ArgNicName, "", "", "The name of the NIC")
	create.AddStringSliceFlag(config.ArgNicIps, "", []string{""}, "IPs assigned to the NIC. This can be a collection")
	create.AddBoolFlag(config.ArgNicDhcp, "", config.DefaultNicDhcp, "Set to false if you wish to disable DHCP on the NIC")
	create.AddIntFlag(config.ArgLanId, "", config.DefaultNicLanId, "The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created")
	create.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcServerIdsNicIdValidate, RunNicUpdate, "update", "Update a NIC",
		`Use this command to update the configuration of a specified NIC. Some restrictions are in place: The primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer. You can also add additional reserved, public IPs to the NIC.

The user can specify and assign private IPs manually. Valid IP addresses for private networks are 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16.

The value for firewallActive can be toggled between true and false to enable or disable the firewall. When the firewall is enabled, incoming traffic is filtered by all the firewall rules configured on the NIC. When the firewall is disabled, then all incoming traffic is routed directly to the NIC and any configured firewall rules are ignored.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* NIC Id`, updateNicExample, true)
	update.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	update.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgNicName, "", "", "The name of the NIC")
	update.AddIntFlag(config.ArgLanId, "", config.DefaultNicLanId, "The LAN ID the NIC sits on")
	update.Command.RegisterFlagCompletionFunc(config.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLansIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgNicDhcp, "", config.DefaultNicDhcp, "Boolean value that indicates if the NIC is using DHCP (true) or not (false)")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcServerIdsNicIdValidate, RunNicDelete, "delete", "Delete a NIC",
		`This command deletes a specified NIC.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id`, deleteNicExample, true)
	deleteCmd.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be deleted [seconds]")

	/*
		Attach Command
	*/
	attachNic := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcIdNicLoadbalancerIdsValidate, RunNicAttach, "attach", "Attach a NIC to a Load Balancer",
		`Use this command to attach a specified NIC to a Load Balancer on your account.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id

The sub-commands of `+"`"+`ionosctl nic attach`+"`"+` allow you to retrieve information about attached NICs or about a specified attached NIC.`, attachNicExample, true)
	attachNic.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	attachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	attachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	attachNic.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to attach to a Load Balancer")
	attachNic.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be attached to a Load Balancer [seconds]")

	/*
		Attach List Command
	*/
	listAttached := builder.NewCommand(context.TODO(), attachNic, PreRunAttachGlobalDcIdLoadbalancerIdValidate, RunNicsAttachList, "list", "List attached NICs from a Load Balancer",
		"Use this command to get a list of attached NICs to a Load Balancer from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Load Balancer Id",
		attachListNicExample, true)
	listAttached.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	listAttached.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Attach Get Command
	*/
	getAttached := builder.NewCommand(context.TODO(), attachNic, PreRunAttachGlobalDcIdNicLoadbalancerIdsValidate, RunNicAttachGet, "get", "Get an attached NIC to a Load Balancer",
		"Use this command to retrieve information about an attached NIC.\n\nRequired values to run the command:\n\n* Data Center Id\n* Load Balancer Id\n* NIC Id",
		attachGetNicExample, true)
	getAttached.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	getAttached.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	getAttached.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	getAttached.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(os.Stderr, nicCmd.Command.Name(), getAttached.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Command
	*/
	detachNic := builder.NewCommand(context.TODO(), nicCmd, PreRunGlobalDcIdNicLoadbalancerIdsValidate, RunNicDetach, "detach", "Detach a NIC from a Load Balancer",
		`Use this command to detach a NIC from a Load Balancer.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`, detachNicExample, true)
	detachNic.AddStringFlag(config.ArgNicId, "", "", config.RequiredFlagNicId)
	detachNic.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedNicsIds(os.Stderr, nicCmd.Command.Name(), detachNic.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddStringFlag(config.ArgLoadBalancerId, "", "", config.RequiredFlagLoadBalancerId)
	detachNic.Command.RegisterFlagCompletionFunc(config.ArgLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLoadbalancersIds(os.Stderr, nicCmd.Command.Name()), cobra.ShellCompDirectiveNoFileComp
	})
	detachNic.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for NIC to detach from a Load Balancer")
	detachNic.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for NIC to be detached from a Load Balancer [seconds]")

	return nicCmd
}

func PreRunGlobalDcServerIdsValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId)
}

func PreRunGlobalDcServerIdsNicIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgNicId)
	if err != nil {
		return err
	}
	return nil
}

func PreRunGlobalDcIdNicLoadbalancerIdsValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgNicId, config.ArgLoadBalancerId)
	if err != nil {
		return err
	}
	return nil
}

func RunNicList(c *builder.CommandConfig) error {
	nics, _, err := c.Nics().List(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)))
	if err != nil {
		return err
	}
	ss := getNics(nics)
	return c.Printer.Print(printer.Result{
		OutputJSON: nics,
		KeyValue:   getNicsKVMaps(ss),
		Columns:    getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunNicGet(c *builder.CommandConfig) error {
	nic, _, err := c.Nics().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: nic,
		KeyValue:   getNicsKVMaps([]resources.Nic{*nic}),
		Columns:    getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})

}

func RunNicCreate(c *builder.CommandConfig) error {
	nic, resp, err := c.Nics().Create(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicName)),
		viper.GetStringSlice(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicIps)),
		viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicDhcp)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  nic,
		KeyValue:    getNicsKVMaps([]resources.Nic{*nic}),
		Columns:     getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "nic",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunNicUpdate(c *builder.CommandConfig) error {
	input := resources.NicProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicName)) {
		input.NicProperties.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicDhcp)) {
		input.NicProperties.SetDhcp(viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicDhcp)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)) {
		input.NicProperties.SetLan(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgLanId)))
	}
	nicUpd, resp, err := c.Nics().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
		input,
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  nicUpd,
		KeyValue:    getNicsKVMaps([]resources.Nic{*nicUpd}),
		Columns:     getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "nic",
		Verb:        "update",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunNicDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nic")
	if err != nil {
		return err
	}
	resp, err := c.Nics().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "nic",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func PreRunAttachGlobalDcIdLoadbalancerIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags("nic", config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLoadBalancerId)
	if err != nil {
		return err
	}
	return nil
}

func PreRunAttachGlobalDcIdNicLoadbalancerIdsValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags("nic", config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgNicId, config.ArgLoadBalancerId)
	if err != nil {
		return err
	}
	return nil
}

func RunNicAttach(c *builder.CommandConfig) error {
	attachedNic, resp, err := c.Nics().AttachToLoadBalancer(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  attachedNic,
		KeyValue:    getNicsKVMaps([]resources.Nic{*attachedNic}),
		Columns:     getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "nic",
		Verb:        "attach",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunNicsAttachList(c *builder.CommandConfig) error {
	attachedNics, _, err := c.Nics().ListAttachedToLoadBalancer(
		viper.GetString(builder.GetGlobalFlagName("nic", config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	vs := getAttachedNics(attachedNics)
	return c.Printer.Print(printer.Result{
		OutputJSON: attachedNics,
		KeyValue:   getNicsKVMaps(vs),
		Columns:    getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunNicAttachGet(c *builder.CommandConfig) error {
	nic, _, err := c.Nics().GetAttachedToLoadBalancer(
		viper.GetString(builder.GetGlobalFlagName("nic", config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: nic,
		KeyValue:   getNicsKVMaps([]resources.Nic{*nic}),
		Columns:    getNicsCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunNicDetach(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "detach nic")
	if err != nil {
		return err
	}
	resp, err := c.Nics().DetachFromLoadBalancer(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLoadBalancerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgNicId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "nic",
		Verb:        "detach",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultNicCols = []string{"NicId", "Name", "Dhcp", "LanId", "Ips"}

type NicPrint struct {
	NicId          string   `json:"NicId,omitempty"`
	Name           string   `json:"Name,omitempty"`
	Dhcp           bool     `json:"Dhcp,omitempty"`
	LanId          int32    `json:"LanId,omitempty"`
	Ips            []string `json:"Ips,omitempty"`
	FirewallActive bool     `json:"FirewallActive,omitempty"`
	Mac            string   `json:"Mac,omitempty"`
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
		ns = append(ns, resources.Nic{s})
	}
	return ns
}

func getAttachedNics(nics resources.BalancedNics) []resources.Nic {
	ns := make([]resources.Nic, 0)
	for _, s := range *nics.BalancedNics.Items {
		ns = append(ns, resources.Nic{s})
	}
	return ns
}

func getNicsKVMaps(ns []resources.Nic) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ns))
	for _, n := range ns {
		properties := n.GetProperties()
		var nicprint NicPrint
		if id, ok := n.GetIdOk(); ok && id != nil {
			nicprint.NicId = *id
		}
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
		o := structs.Map(nicprint)
		out = append(out, o)
	}
	return out
}

func getNicsIds(outErr io.Writer, parentCmdName string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)

	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)

	nicSvc := resources.NewNicService(clientSvc.Get(), context.TODO())
	nics, _, err := nicSvc.List(
		viper.GetString(builder.GetGlobalFlagName(parentCmdName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(parentCmdName, config.ArgServerId)),
	)
	clierror.CheckError(err, outErr)

	nicsIds := make([]string, 0)
	if nics.Nics.Items != nil {
		for _, v := range *nics.Nics.Items {
			nicsIds = append(nicsIds, *v.GetId())
		}
	} else {
		return nil
	}
	return nicsIds
}

func getAttachedNicsIds(outErr io.Writer, parentCmdName, nameCmd string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)

	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)

	nicSvc := resources.NewNicService(clientSvc.Get(), context.TODO())
	nics, _, err := nicSvc.ListAttachedToLoadBalancer(
		viper.GetString(builder.GetGlobalFlagName(parentCmdName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(parentCmdName, nameCmd, config.ArgLoadBalancerId)),
	)
	clierror.CheckError(err, outErr)

	attachedNicsIds := make([]string, 0)
	if nics.BalancedNics.Items != nil {
		for _, v := range *nics.BalancedNics.Items {
			attachedNicsIds = append(attachedNicsIds, *v.GetId())
		}
	} else {
		return nil
	}
	return attachedNicsIds
}
