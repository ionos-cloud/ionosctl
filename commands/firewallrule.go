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
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func firewallrule() *builder.Command {
	ctx := context.TODO()
	firewallRuleCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "firewallrule",
			Aliases:          []string{"rule"},
			Short:            "Firewall Rule Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl firewallrule` + "`" + ` allow you to create, list, get, update, delete Firewall Rules.`,
			TraverseChildren: true,
		},
	}
	globalFlags := firewallRuleCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgNicId), globalFlags.Lookup(config.ArgNicId))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultFirewallRuleCols, "Columns to be printed in the standard output. Example: --cols \"ResourceId,Name\"")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, firewallRuleCmd, PreRunGlobalDcServerNicIdsValidate, RunFirewallRuleList, "list", "List Firewall Rules",
		"Use this command to get a list of Firewall Rules from a specified NIC from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n*Nic Id",
		listFirewallRuleExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, firewallRuleCmd, PreRunGlobalDcServerNicIdsFRuleIdValidate, RunFirewallRuleGet, "get", "Get a Firewall Rule",
		"Use this command to retrieve information of a specified Firewall Rule.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n*Nic Id\n* FirewallRule Id",
		getFirewallRuleExample, true)
	get.AddStringFlag(config.ArgFirewallRuleId, "", "", config.RequiredFlagFirewallRuleId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFirewallRulesIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgServerId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, firewallRuleCmd, PreRunGlobalDcIdValidate, RunFirewallRuleCreate, "create", "Create a Firewall Rule",
		`Use this command to create a new Firewall Rule. Please Note: the Firewall Rule Protocol can only be set when creating a new Firewall Rule.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id 
* Firewall Rule Protocol`, createFirewallRuleExample, true)
	create.AddStringFlag(config.ArgFirewallRuleName, "", "", "The name for the Firewall Rule")
	create.AddStringFlag(config.ArgFirewallRuleProtocol, "", "", "The Protocol for Firewall Rule: TCP, UDP, ICMP, ANY "+config.RequiredFlag)
	create.AddStringFlag(config.ArgFirewallRuleSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses.")
	create.AddStringFlag(config.ArgFirewallRuleSourceIp, "", "", "Only traffic originating from the respective IPv4 address is allowed. Unset option allows all source IPs.")
	create.AddStringFlag(config.ArgFirewallRuleTargetIp, "", "", "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Unset option allows all target IPs.")
	create.AddIntFlag(config.ArgFirewallRuleIcmpType, "", 0, "Define the allowed type (from 0 to 254) if the protocol ICMP is chosen. Unset option allows all types.")
	create.AddIntFlag(config.ArgFirewallRuleIcmpCode, "", 0, "Define the allowed code (from 0 to 254) if protocol ICMP is chosen. Unset option allows all codes.")
	create.AddIntFlag(config.ArgFirewallRulePortRangeStart, "", 1, "Define the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.")
	create.AddIntFlag(config.ArgFirewallRulePortRangeStop, "", 1, "Define the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.")
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Firewall Rule to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Firewall Rule to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, firewallRuleCmd, PreRunGlobalDcServerNicIdsFRuleIdValidate, RunFirewallRuleUpdate, "update", "Update a FirewallRule",
		`Use this command to update a specified Firewall Rule.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id`, updateFirewallRuleExample, true)
	update.AddStringFlag(config.ArgFirewallRuleName, "", "", "The name for the Firewall Rule")
	update.AddStringFlag(config.ArgFirewallRuleSourceMac, "", "", "Only traffic originating from the respective MAC address is allowed. Valid format: aa:bb:cc:dd:ee:ff. Unset option allows all source MAC addresses.")
	update.AddStringFlag(config.ArgFirewallRuleSourceIp, "", "", "Only traffic originating from the respective IPv4 address is allowed. Unset option allows all source IPs.")
	update.AddStringFlag(config.ArgFirewallRuleTargetIp, "", "", "In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed. Unset option allows all target IPs.")
	update.AddIntFlag(config.ArgFirewallRuleIcmpType, "", 0, "Redefine the allowed type (from 0 to 254) if the protocol ICMP is chosen. Unset option allows all types.")
	update.AddIntFlag(config.ArgFirewallRuleIcmpCode, "", 0, "Redefine the allowed code (from 0 to 254) if protocol ICMP is chosen. Unset option allows all codes.")
	update.AddIntFlag(config.ArgFirewallRulePortRangeStart, "", 1, "Redefine the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.")
	update.AddIntFlag(config.ArgFirewallRulePortRangeStop, "", 1, "Redefine the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen. Leave portRangeStart and portRangeEnd unset to allow all ports.")
	update.AddStringFlag(config.ArgFirewallRuleId, "", "", config.RequiredFlagFirewallRuleId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFirewallRulesIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgServerId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Firewall Rule to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Firewall Rule to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, firewallRuleCmd, PreRunGlobalDcServerNicIdsFRuleIdValidate, RunFirewallRuleDelete, "delete", "Delete a FirewallRule",
		`Use this command to delete a specified Firewall Rule from a Virtual Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id`, deleteFirewallRuleExample, true)
	deleteCmd.AddStringFlag(config.ArgFirewallRuleId, "", "", config.RequiredFlagFirewallRuleId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgFirewallRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFirewallRulesIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgServerId)),
			viper.GetString(builder.GetGlobalFlagName(firewallRuleCmd.Command.Name(), config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Firewall Rule to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Firewall Rule to be deleted [seconds]")

	return firewallRuleCmd
}

func PreRunGlobalDcServerNicIdsValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId)
}

func PreRunGlobalDcServerNicIdsFRuleIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgFirewallRuleId)
	if err != nil {
		return err
	}
	return nil
}

func RunFirewallRuleList(c *builder.CommandConfig) error {
	firewallRules, _, err := c.FirewallRules().List(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(nil, c, getFirewallRules(firewallRules)))
}

func RunFirewallRuleGet(c *builder.CommandConfig) error {
	firewallRule, _, err := c.FirewallRules().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgNicId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(nil, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleCreate(c *builder.CommandConfig) error {
	properties := getFirewallRulePropertiesSet(c)
	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	firewallRule, resp, err := c.FirewallRules().Create(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgNicId)),
		input,
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(resp, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleUpdate(c *builder.CommandConfig) error {
	firewallRule, resp, err := c.FirewallRules().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgNicId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleId)),
		getFirewallRulePropertiesSet(c),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(resp, c, getFirewallRule(firewallRule)))
}

func RunFirewallRuleDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete firewall rule")
	if err != nil {
		return err
	}
	resp, err := c.FirewallRules().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgNicId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getFirewallRulePrint(resp, c, nil))
}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *builder.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleName)) {
		properties.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleProtocol)) {
		properties.SetProtocol(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleProtocol)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleSourceIp)) {
		properties.SetSourceIp(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleSourceIp)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleSourceMac)) {
		properties.SetSourceMac(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleSourceMac)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleTargetIp)) {
		properties.SetTargetIp(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleTargetIp)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleIcmpCode)) {
		properties.SetIcmpCode(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleIcmpCode)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleIcmpType)) {
		properties.SetIcmpType(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRuleIcmpType)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRulePortRangeStart)) {
		properties.SetPortRangeStart(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRulePortRangeStart)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRulePortRangeStop)) {
		properties.SetPortRangeEnd(viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgFirewallRulePortRangeStop)))
	}
	return properties
}

// Output Printing

var defaultFirewallRuleCols = []string{"FirewallRuleId", "Name", "Protocol", "PortRangeStart", "PortRangeEnd", "State"}

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

func getFirewallRulePrint(resp *resources.Response, c *builder.CommandConfig, s []resources.FirewallRule) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getFirewallRulesKVMaps(s)
			r.Columns = getFirewallRulesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getFirewallRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultFirewallRuleCols
	}

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
	var FirewallRuleCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			FirewallRuleCols = append(FirewallRuleCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return FirewallRuleCols
}

func getFirewallRules(firewallRules resources.FirewallRules) []resources.FirewallRule {
	ls := make([]resources.FirewallRule, 0)
	for _, s := range *firewallRules.Items {
		ls = append(ls, resources.FirewallRule{FirewallRule: s})
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
	for _, l := range ls {
		o := getFirewallRuleKVMap(l)
		out = append(out, o)
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

func getFirewallRulesIds(outErr io.Writer, datacenterId, serverId, nicId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	firewallRuleSvc := resources.NewFirewallRuleService(clientSvc.Get(), context.TODO())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId)
	clierror.CheckError(err, outErr)
	firewallRulesIds := make([]string, 0)
	if items, ok := firewallRules.FirewallRules.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				firewallRulesIds = append(firewallRulesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return firewallRulesIds
}
