package commands

import (
	"context"
	"errors"
	"fmt"
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

func networkloadbalancerRule() *core.Command {
	ctx := context.TODO()
	nlbRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Network Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl nlb rule` allow you to create, list, get, update, delete Network Load Balancer Forwarding Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultForwardingRuleCols, utils.ColsMessage(allForwardingRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nlbRuleCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nlbRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "rule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancer Forwarding Rules",
		LongDesc:   "Use this command to list Network Load Balancer Forwarding Rules from a specified Network Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id",
		Example:    listNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace:  "networkloadbalancer",
		Resource:   "rule",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Network Load Balancer Forwarding Rule",
		LongDesc:   "Use this command to get information about a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id\n* Forwarding Rule Id",
		Example:    getNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagForwardingRuleId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a Network Load Balancer Forwarding Rule in a specified Network Load Balancer. You can also set Health Check Settings for Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Listener Ip
* Listener Port`,
		Example:    createNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunNetworkLoadBalancerForwardingRuleCreate,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Unnamed Forwarding Rule", "The name for the Forwarding Rule")
	create.AddStringFlag(config.ArgListenerIp, "", "", "Listening IP "+config.RequiredFlag)
	create.AddStringFlag(config.ArgListenerPort, "", "", "Listening port number. Range: 1 to 65535 "+config.RequiredFlag)
	create.AddIntFlag(config.ArgRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	create.AddIntFlag(config.ArgClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	create.AddIntFlag(config.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	create.AddIntFlag(config.ArgTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	create.AddStringFlag(config.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to update a specified Network Load Balancer Forwarding Rule from a Network Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example:    updateNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagForwardingRuleId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(update.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the Forwarding Rule")
	update.AddStringFlag(config.ArgListenerIp, "", "", "Listening IP "+config.RequiredFlag)
	update.AddStringFlag(config.ArgListenerPort, "", "", "Listening port number. Range: 1 to 65535 "+config.RequiredFlag)
	update.AddStringFlag(config.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(config.ArgRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	update.AddIntFlag(config.ArgClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	update.AddIntFlag(config.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	update.AddIntFlag(config.ArgTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, nlbRuleCmd, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to delete a specified Network Load Balancer Forwarding Rule from a Network Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example:    deleteNetworkLoadBalancerForwardingRuleExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgNetworkLoadBalancerId, "", "", config.RequiredFlagNetworkLoadBalancerId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgRuleId, config.ArgIdShort, "", config.RequiredFlagForwardingRuleId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Forwarding Rule deletion [seconds]")

	nlbRuleCmd.AddCommand(nlbRuleTarget())

	return nlbRuleCmd
}

func PreRunNetworkLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNetworkLoadBalancerId, config.ArgListenerIp, config.ArgListenerPort)
}

func PreRunDcNetworkLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgNetworkLoadBalancerId, config.ArgRuleId)
}

func RunNetworkLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	nlbForwardingRules, _, err := c.NetworkLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(nil, c, getForwardingRules(nlbForwardingRules)))
}

func RunNetworkLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	ng, _, err := c.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(nil, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	proper := getForwardingRulePropertiesSet(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(string(ionoscloud.TCP))
	}
	if !proper.HasAlgorithm() {
		proper.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, config.ArgAlgorithm)))
	}
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		proper.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}
	ng, resp, err := c.NetworkLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		resources.NetworkLoadBalancerForwardingRule{
			NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
				Properties: &proper.NetworkLoadBalancerForwardingRuleProperties,
			},
		},
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(resp, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	input := getForwardingRulePropertiesSet(c)
	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		input.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}
	ng, resp, err := c.NetworkLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
		input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(resp, c, []resources.NetworkLoadBalancerForwardingRule{*ng}))
}

func RunNetworkLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete network load balancer forwarding rule"); err != nil {
		return err
	}
	resp, err := c.NetworkLoadBalancers().DeleteForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getForwardingRulePrint(resp, c, nil))
}

func getForwardingRulePropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleProperties {
	input := ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgAlgorithm)) {
		input.SetAlgorithm(strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgAlgorithm))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgListenerIp)) {
		input.SetListenerIp(viper.GetString(core.GetFlagName(c.NS, config.ArgListenerIp)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgListenerPort)) {
		input.SetListenerPort(viper.GetInt32(core.GetFlagName(c.NS, config.ArgListenerPort)))
	}
	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: input,
	}
}

func getForwardingRuleHealthCheckPropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleHealthCheck {
	inputHealth := ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgRetries)) {
		inputHealth.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, config.ArgRetries)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgClientTimeout)) {
		inputHealth.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, config.ArgClientTimeout)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgConnectionTimeout)) {
		inputHealth.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, config.ArgConnectionTimeout)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgTargetTimeout)) {
		inputHealth.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, config.ArgTargetTimeout)))
	}
	return &resources.NetworkLoadBalancerForwardingRuleHealthCheck{
		NetworkLoadBalancerForwardingRuleHealthCheck: inputHealth,
	}
}

// Output Printing

var (
	defaultForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State"}
	allForwardingRuleCols     = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State",
		"ClientTimeout", "ConnectTimeout", "TargetTimeout", "Retries"}
)

type ForwardingRulePrint struct {
	ForwardingRuleId string `json:"ForwardingRuleId,omitempty"`
	Name             string `json:"Name,omitempty"`
	Algorithm        string `json:"Algorithm,omitempty"`
	Protocol         string `json:"Protocol,omitempty"`
	ListenerIp       string `json:"ListenerIp,omitempty"`
	ListenerPort     int32  `json:"ListenerPort,omitempty"`
	ClientTimeout    string `json:"ClientTimeout,omitempty"`
	ConnectTimeout   string `json:"ConnectTimeout,omitempty"`
	TargetTimeout    string `json:"TargetTimeout,omitempty"`
	Retries          int32  `json:"Retries,omitempty"`
	State            string `json:"State,omitempty"`
}

func getForwardingRulePrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NetworkLoadBalancerForwardingRule) printer.Result {
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
			r.KeyValue = getForwardingRulesKVMaps(ss)
			r.Columns = getForwardingRulesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getForwardingRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultForwardingRuleCols
	}

	columnsMap := map[string]string{
		"ForwardingRuleId": "ForwardingRuleId",
		"Name":             "Name",
		"Algorithm":        "Algorithm",
		"Protocol":         "Protocol",
		"ListenerIp":       "ListenerIp",
		"ListenerPort":     "ListenerPort",
		"ClientTimeout":    "ClientTimeout",
		"ConnectTimeout":   "ConnectTimeout",
		"TargetTimeout":    "TargetTimeout",
		"Retries":          "Retries",
		"State":            "State",
	}
	var forwardingRuleCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			forwardingRuleCols = append(forwardingRuleCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return forwardingRuleCols
}

func getForwardingRules(forwardingrules resources.NetworkLoadBalancerForwardingRules) []resources.NetworkLoadBalancerForwardingRule {
	ss := make([]resources.NetworkLoadBalancerForwardingRule, 0)
	for _, s := range *forwardingrules.Items {
		ss = append(ss, resources.NetworkLoadBalancerForwardingRule{NetworkLoadBalancerForwardingRule: s})
	}
	return ss
}

func getForwardingRulesKVMaps(ss []resources.NetworkLoadBalancerForwardingRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		var forwardingRulePrint ForwardingRulePrint
		if id, ok := s.GetIdOk(); ok && id != nil {
			forwardingRulePrint.ForwardingRuleId = *id
		}
		if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				forwardingRulePrint.Name = *name
			}
			if listenerIp, ok := properties.GetListenerIpOk(); ok && listenerIp != nil {
				forwardingRulePrint.ListenerIp = *listenerIp
			}
			if listenerPort, ok := properties.GetListenerPortOk(); ok && listenerPort != nil {
				forwardingRulePrint.ListenerPort = *listenerPort
			}
			if protocol, ok := properties.GetProtocolOk(); ok && protocol != nil {
				forwardingRulePrint.Protocol = *protocol
			}
			if algorithm, ok := properties.GetAlgorithmOk(); ok && algorithm != nil {
				forwardingRulePrint.Algorithm = *algorithm
			}
			if health, ok := properties.GetHealthCheckOk(); ok && health != nil {
				if clientTimeout, ok := health.GetClientTimeoutOk(); ok && clientTimeout != nil {
					forwardingRulePrint.ClientTimeout = fmt.Sprintf("%vms", *clientTimeout)
				}
				if connectTimeout, ok := health.GetConnectTimeoutOk(); ok && connectTimeout != nil {
					forwardingRulePrint.ConnectTimeout = fmt.Sprintf("%vms", *connectTimeout)
				}
				if targetTimeout, ok := health.GetTargetTimeoutOk(); ok && targetTimeout != nil {
					forwardingRulePrint.TargetTimeout = fmt.Sprintf("%vms", *targetTimeout)
				}
				if retries, ok := health.GetRetriesOk(); ok && retries != nil {
					forwardingRulePrint.Retries = *retries
				}
			}
		}
		if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				forwardingRulePrint.State = *state
			}
		}
		o := structs.Map(forwardingRulePrint)
		out = append(out, o)
	}
	return out
}

func getForwardingRulesIds(outErr io.Writer, datacenterId, nlbId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	nlbSvc := resources.NewNetworkLoadBalancerService(clientSvc.Get(), context.TODO())
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId)
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.NetworkLoadBalancerForwardingRules.GetItemsOk(); ok && items != nil {
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
