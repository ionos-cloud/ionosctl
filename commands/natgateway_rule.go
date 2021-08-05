package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func natgatewayRule() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultNatGatewayRuleCols, utils.ColsMessage(allNatGatewayRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(natgatewayRuleCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = natgatewayRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		LongDesc:   "Use this command to list NAT Gateway Rules from a specified NAT Gateway.\n\nRequired values to run command:\n\n* Data Center Id\n* NAT Gateway Id",
		Example:    listNatGatewayRuleExample,
		PreCmdRun:  PreRunDcNatGatewayIds,
		CmdRun:     RunNatGatewayRuleList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagRuleId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})

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
* Name
* Public IP
* Source Subnet`,
		Example:    createNatGatewayRuleExample,
		PreCmdRun:  PreRunNatGatewayRuleCreate,
		CmdRun:     RunNatGatewayRuleCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the NAT Gateway Rule "+config.RequiredFlag)
	create.AddStringFlag(config.ArgProtocol, config.ArgProtocolShort, string(ionoscloud.ALL), "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgIp, "", "", "Public IP address of the NAT Gateway Rule "+config.RequiredFlag)
	create.AddStringFlag(config.ArgSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule "+config.RequiredFlag)
	create.AddStringFlag(config.ArgTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	create.AddIntFlag(config.ArgPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	create.AddIntFlag(config.ArgPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule creation [seconds]")

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
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagRuleId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, config.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the NAT Gateway Rule")
	update.AddStringFlag(config.ArgProtocol, config.ArgProtocolShort, "", "Protocol of the NAT Gateway Rule. If protocol is 'ICMP' then targetPortRange start and end cannot be set")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{string(ionoscloud.TCP), string(ionoscloud.UDP), string(ionoscloud.ICMP), string(ionoscloud.ALL)}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgIp, "", "", "Public IP address of the NAT Gateway Rule")
	update.AddStringFlag(config.ArgSourceSubnet, "", "", "Source subnet of the NAT Gateway Rule")
	update.AddStringFlag(config.ArgTargetSubnet, "", "", "Target subnet or destination subnet of the NAT Gateway Rule")
	update.AddIntFlag(config.ArgPortRangeStart, "", 1, "Target port range start associated with the NAT Gateway Rule")
	update.AddIntFlag(config.ArgPortRangeEnd, "", 1, "Target port range end associated with the NAT Gateway Rule")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule update [seconds]")

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
		PreCmdRun:  PreRunDcNatGatewayRuleIds,
		CmdRun:     RunNatGatewayRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgNatGatewayId, "", "", config.RequiredFlagNatGatewayId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNatGatewayId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewaysIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagRuleId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNatGatewayRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgNatGatewayId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for NAT Gateway Rule deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for NAT Gateway Rule deletion [seconds]")

	return natgatewayRuleCmd
}

func PreRunNatGatewayRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNatGatewayId, config.ArgName, config.ArgIp, config.ArgSourceSubnet)
}

func PreRunDcNatGatewayRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNatGatewayId, config.ArgRuleId)
}

func RunNatGatewayRuleList(c *core.CommandConfig) error {
	natgatewayRules, _, err := c.NatGateways().ListRules(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(nil, c, getNatGatewayRules(natgatewayRules)))
}

func RunNatGatewayRuleGet(c *core.CommandConfig) error {
	ng, _, err := c.NatGateways().GetRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(nil, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleCreate(c *core.CommandConfig) error {
	proper := getNewNatGatewayRuleInfo(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(ionoscloud.NatGatewayRuleProtocol(viper.GetString(core.GetFlagName(c.NS, config.ArgProtocol))))
	}
	ng, resp, err := c.NatGateways().CreateRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
		resources.NatGatewayRule{
			NatGatewayRule: ionoscloud.NatGatewayRule{
				Properties: &proper.NatGatewayRuleProperties,
			},
		},
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(resp, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleUpdate(c *core.CommandConfig) error {
	input := getNewNatGatewayRuleInfo(c)
	ng, resp, err := c.NatGateways().UpdateRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
		*input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(resp, c, []resources.NatGatewayRule{*ng}))
}

func RunNatGatewayRuleDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete nat gateway rule"); err != nil {
		return err
	}
	resp, err := c.NatGateways().DeleteRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getNatGatewayRulePrint(resp, c, nil))
}

func getNewNatGatewayRuleInfo(c *core.CommandConfig) *resources.NatGatewayRuleProperties {
	input := ionoscloud.NatGatewayRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIp)) {
		input.SetPublicIp(viper.GetString(core.GetFlagName(c.NS, config.ArgIp)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgProtocol)) {
		protocol := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgProtocol)))
		input.SetProtocol(ionoscloud.NatGatewayRuleProtocol(protocol))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSourceSubnet)) {
		input.SetSourceSubnet(viper.GetString(core.GetFlagName(c.NS, config.ArgSourceSubnet)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgTargetSubnet)) {
		input.SetTargetSubnet(viper.GetString(core.GetFlagName(c.NS, config.ArgTargetSubnet)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPortRangeStart)) &&
		viper.IsSet(core.GetFlagName(c.NS, config.ArgPortRangeEnd)) {
		inputPortRange := ionoscloud.TargetPortRange{}
		inputPortRange.SetStart(viper.GetInt32(core.GetFlagName(c.NS, config.ArgPortRangeStart)))
		inputPortRange.SetEnd(viper.GetInt32(core.GetFlagName(c.NS, config.ArgPortRangeEnd)))
		input.SetTargetPortRange(inputPortRange)
	}
	return &resources.NatGatewayRuleProperties{
		NatGatewayRuleProperties: input,
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
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getNatGatewayRulesKVMaps(ss)
			r.Columns = getNatGatewayRulesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
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
	ss := make([]resources.NatGatewayRule, 0)
	for _, s := range *natgatewayRules.Items {
		ss = append(ss, resources.NatGatewayRule{NatGatewayRule: s})
	}
	return ss
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

func getNatGatewayRulesIds(outErr io.Writer, datacenterId, natgatewayId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	natgatewaySvc := resources.NewNatGatewayService(clientSvc.Get(), context.TODO())
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := natgateways.NatGatewayRules.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
