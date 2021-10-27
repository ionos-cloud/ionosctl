package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AlbRuleHttpRuleCmd() *core.Command {
	ctx := context.TODO()
	albRuleHttpRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "httprule",
			Aliases:          []string{"http"},
			Short:            "Application Load Balancer Forwarding Rule Http Rule Operations",
			Long:             "The sub-commands of `ionosctl alb rule httprule` allow you to add, list, update, remove Application Load Balancer Forwarding Rule Http Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := albRuleHttpRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultAlbRuleHttpRuleCols, printer.ColsMessage(defaultAlbRuleHttpRuleCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(albRuleHttpRuleCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = albRuleHttpRuleCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAlbRuleHttpRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, albRuleHttpRuleCmd, core.CommandBuilder{
		Namespace:  "forwardingrule",
		Resource:   "httprule",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Application Load Balancer Forwarding Rule Http Rules",
		LongDesc:   "Use this command to list Http Rules of a Application Load Balancer Forwarding Rule.\n\nRequired values to run command:\n\n* Data Center Id\n* Application Load Balancer Id\n* Forwarding Rule Id",
		Example:    "",
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunAlbRuleHttpRuleList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, albRuleHttpRuleCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "httprule",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Http Rule to Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to add a Http Rule in a specified Application Load Balancer Forwarding Rule. 

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name
* Http Rule Type`,
		Example:    "",
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRule,
		CmdRun:     RunAlbRuleHttpRuleAdd,
		InitClient: true,
	})
	add.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Http Rule", "A name of that Application Load Balancer Http Rule", core.RequiredFlagOption())
	add.AddStringFlag(cloudapiv6.ArgType, "", "", "Type of the Http Rule", core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgTargetGroupId, "", "", "The Id of the Target Group; mandatory for FORWARD action")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddBoolFlag(cloudapiv6.ArgDropQuery, "", false, "Default is false; must be true for REDIRECT action")
	add.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "www.ionos.com", "The location for redirecting; mandatory for REDIRECT action")
	add.AddStringFlag(cloudapiv6.ArgStatusCode, "", "301", "On REDIRECT action it can take the value 301, 302, 303, 307, 308; on STATIC action it is between 200 and 599")
	add.AddStringFlag(cloudapiv6.ArgResponse, "", "", "The response message of the request; mandatory for STATIC action")
	add.AddStringFlag(cloudapiv6.ArgContentType, "", "application/json", "")
	add.AddStringFlag(cloudapiv6.ArgCondition, "", "", "Condition of the Http Rule condition")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCondition, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgConditionType, "", "HEADER", "Type of the Http Rule condition")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgConditionType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "Specifies whether the condition is negated or not; default: false")
	add.AddStringFlag(cloudapiv6.ArgConditionKey, "", "forward-at", "")
	add.AddStringFlag(cloudapiv6.ArgConditionValue, "", "Friday", "")
	add.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule creation to be executed")
	add.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule creation [seconds]")

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, albRuleHttpRuleCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "httprule",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Http Rule from a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to remove a specified Http Rule from Application Load Balancer Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name
* Http Rule Type`,
		Example:    "",
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRule,
		CmdRun:     RunAlbRuleHttpRuleRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	// TODO: check if the name is unique, that means type is not needed
	removeCmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Http Rule", "A name of that Application Load Balancer Http Rule", core.RequiredFlagOption())
	removeCmd.AddStringFlag(cloudapiv6.ArgType, "", "", "Type of the Http Rule", core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule deletion to be executed")
	removeCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule deletion [seconds]")

	return albRuleHttpRuleCmd
}

func PreRunApplicationLoadBalancerRuleHttpRule(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId,
		cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgName, cloudapiv6.ArgType)
}

func RunAlbRuleHttpRuleList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting HttpRules from ForwardingRule with ID: %v from ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
		if httpRulesOk, ok := properties.GetHttpRulesOk(); ok && httpRulesOk != nil {
			return c.Printer.Print(getAlbRuleHttpRulePrint(nil, c, getAlbRuleHttpRules(httpRulesOk)))
		} else {
			return errors.New("error getting rule http rules")
		}
	} else {
		return errors.New("error getting rule properties")
	}
}

func RunAlbRuleHttpRuleAdd(c *core.CommandConfig) error {
	var httpRuleItems []ionoscloud.ApplicationLoadBalancerHttpRule
	c.Printer.Verbose("Getting HttpRules from ForwardingRule with ID: %v from ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	ngOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	if properties, ok := ngOld.GetPropertiesOk(); ok && properties != nil {
		if httpRulesOk, ok := properties.GetHttpRulesOk(); ok && httpRulesOk != nil {
			httpRuleItems = *httpRulesOk
		}
	}
	httpRuleNew := getRuleHttpRuleInfo(c)
	c.Printer.Verbose("Adding the new HttpRule to the existing HttpRules")
	httpRuleItems = append(httpRuleItems, httpRuleNew.ApplicationLoadBalancerHttpRule)
	c.Printer.Verbose("Updating ForwardingRule with the new HttpRules")
	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		&resources.ApplicationLoadBalancerForwardingRuleProperties{
			ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
				HttpRules: &httpRuleItems,
			},
		},
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbRuleHttpRulePrint(resp, c, []resources.ApplicationLoadBalancerHttpRule{httpRuleNew}))
}

func RunAlbRuleHttpRuleRemove(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove forwarding rule http rule"); err != nil {
		return err
	}
	c.Printer.Verbose("Getting HttpRules from ForwardingRule with ID: %v from ApplicationLoadBalancer with ID: %v from Datacenter with ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	frOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Removing the HttpRule from the existing HttpRules")
	proper, err := getRuleHttpRulesRemove(c, frOld)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating ForwardingRule with the new HttpRules")
	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		proper,
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getAlbRuleHttpRulePrint(resp, c, nil))
}

func getRuleHttpRuleInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRule {
	httprule := resources.ApplicationLoadBalancerHttpRule{}
	httprule.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	httprule.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	c.Printer.Verbose("Property Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	httprule.SetTargetGroup(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Printer.Verbose("Property TargetGroup set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	httprule.SetDropQuery(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDropQuery)))
	c.Printer.Verbose("Property DropQuery set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDropQuery)))
	httprule.SetLocation(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
	c.Printer.Verbose("Property Location set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
	httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
	c.Printer.Verbose("Property StatusCode set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
	httprule.SetResponseMessage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	c.Printer.Verbose("Property ResponseMessage set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	httprule.SetContentType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))
	c.Printer.Verbose("Property ContentType set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))
	httpRuleCondition := resources.ApplicationLoadBalancerHttpRuleCondition{}
	httpRuleCondition.SetCondition(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
	c.Printer.Verbose("Property Condition set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
	httpRuleCondition.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	c.Printer.Verbose("Property Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	httpRuleCondition.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	c.Printer.Verbose("Property Negate set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	httpRuleCondition.SetKey(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
	c.Printer.Verbose("Property Key set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
	httpRuleCondition.SetValue(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
	c.Printer.Verbose("Property Value set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
	httprule.SetConditions([]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{httpRuleCondition.ApplicationLoadBalancerHttpRuleCondition})
	c.Printer.Verbose("Setting Conditions to HttpRule")
	return httprule
}

func getRuleHttpRulesRemove(c *core.CommandConfig, frOld *resources.ApplicationLoadBalancerForwardingRule) (*resources.ApplicationLoadBalancerForwardingRuleProperties, error) {
	var (
		foundIp   = false
		foundPort = false
	)
	httpruleItems := make([]ionoscloud.ApplicationLoadBalancerHttpRule, 0)
	if properties, ok := frOld.GetPropertiesOk(); ok && properties != nil {
		if httprules, ok := properties.GetHttpRulesOk(); ok && httprules != nil {
			// Iterate trough all httprules
			for _, httpruleItem := range *httprules {
				removeName := false
				removeType := false
				if nameOk, ok := httpruleItem.GetNameOk(); ok && nameOk != nil {
					if *nameOk == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
						removeName = true
						foundIp = true
					}
				}
				if typeOk, ok := httpruleItem.GetTypeOk(); ok && typeOk != nil {
					if *typeOk == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)) {
						removeType = true
						foundPort = true
					}
				}
				if removeName && removeType {
					continue
				} else {
					httpruleItems = append(httpruleItems, httpruleItem)
				}
			}
		}
	}
	if !foundIp {
		return nil, errors.New("no forwarding rule http rule with the specified IP found")
	}
	if !foundPort {
		return nil, errors.New("no forwarding rule http rule with the specified port found")
	}
	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &httpruleItems,
		},
	}, nil
}

// Output Printing

var defaultAlbRuleHttpRuleCols = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Location", "StatusCode", "ResponseMessage", "ContentType", "Condition"}

type AlbRuleHttpRulePrint struct {
	Name            string   `json:"Name,omitempty"`
	Type            string   `json:"Type,omitempty"`
	TargetGroupId   string   `json:"TargetGroupId,omitempty"`
	DropQuery       bool     `json:"DropQuery,omitempty"`
	Location        string   `json:"Location,omitempty"`
	StatusCode      int32    `json:"StatusCode,omitempty"`
	ResponseMessage string   `json:"ResponseMessage,omitempty"`
	ContentType     string   `json:"ContentType,omitempty"`
	Condition       []string `json:"Condition,omitempty"`
}

func getAlbRuleHttpRulePrint(resp *resources.Response, c *core.CommandConfig, ss []resources.ApplicationLoadBalancerHttpRule) printer.Result {
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
			r.KeyValue = getAlbRuleHttpRulesKVMaps(ss)
			r.Columns = getAlbRuleHttpRulesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getAlbRuleHttpRulesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultAlbRuleHttpRuleCols
	}

	columnsMap := map[string]string{
		"Name":            "Name",
		"Type":            "Type",
		"TargetGroupId":   "TargetGroupId",
		"DropQuery":       "DropQuery",
		"Location":        "Location",
		"StatusCode":      "StatusCode",
		"ResponseMessage": "ResponseMessage",
		"ContentType":     "ContentType",
		"Condition":       "Condition",
	}
	var ruleHttpRuleCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			ruleHttpRuleCols = append(ruleHttpRuleCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return ruleHttpRuleCols
}

func getAlbRuleHttpRules(httprules *[]ionoscloud.ApplicationLoadBalancerHttpRule) []resources.ApplicationLoadBalancerHttpRule {
	ss := make([]resources.ApplicationLoadBalancerHttpRule, 0)
	if httprules != nil {
		for _, s := range *httprules {
			ss = append(ss, resources.ApplicationLoadBalancerHttpRule{
				ApplicationLoadBalancerHttpRule: s,
			})
		}
	}
	return ss
}

func getAlbRuleHttpRulesKVMaps(httprules []resources.ApplicationLoadBalancerHttpRule) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(httprules))
	for _, httprule := range httprules {
		var httpRulePrint AlbRuleHttpRulePrint
		if nameOk, ok := httprule.GetNameOk(); ok && nameOk != nil {
			httpRulePrint.Name = *nameOk
		}
		if typeOk, ok := httprule.GetTypeOk(); ok && typeOk != nil {
			httpRulePrint.Type = *typeOk
		}
		if targetGroupOk, ok := httprule.GetTargetGroupOk(); ok && targetGroupOk != nil {
			httpRulePrint.TargetGroupId = *targetGroupOk
		}
		if dropQueryOk, ok := httprule.GetDropQueryOk(); ok && dropQueryOk != nil {
			httpRulePrint.DropQuery = *dropQueryOk
		}
		if locationOk, ok := httprule.GetLocationOk(); ok && locationOk != nil {
			httpRulePrint.Location = *locationOk
		}
		if statusCodeOk, ok := httprule.GetStatusCodeOk(); ok && statusCodeOk != nil {
			httpRulePrint.StatusCode = *statusCodeOk
		}
		if responseMessageOk, ok := httprule.GetResponseMessageOk(); ok && responseMessageOk != nil {
			httpRulePrint.ResponseMessage = *responseMessageOk
		}
		if contentTypeOk, ok := httprule.GetContentTypeOk(); ok && contentTypeOk != nil {
			httpRulePrint.ContentType = *contentTypeOk
		}
		if conditionsOk, ok := httprule.GetConditionsOk(); ok && conditionsOk != nil {
			conditions := make([]string, 0)
			for _, conditionOk := range *conditionsOk {
				var condition string
				if getConditionOk, ok := conditionOk.GetConditionOk(); ok && getConditionOk != nil {
					condition = fmt.Sprintf("Condition: %s", *getConditionOk)
				}
				if getTypeOk, ok := conditionOk.GetTypeOk(); ok && getTypeOk != nil {
					condition = fmt.Sprintf("%s Type: %s", condition, *getTypeOk)
				}
				if getNegateOk, ok := conditionOk.GetNegateOk(); ok && getNegateOk != nil {
					condition = fmt.Sprintf("%s Negate: %v", condition, *getNegateOk)
				}
				if getKeyOk, ok := conditionOk.GetKeyOk(); ok && getKeyOk != nil {
					condition = fmt.Sprintf("%s Key: %s", condition, *getKeyOk)
				}
				if getValueOk, ok := conditionOk.GetValueOk(); ok && getValueOk != nil {
					condition = fmt.Sprintf("%s Value: %s", condition, *getValueOk)
				}
				conditions = append(conditions, condition)
			}
			httpRulePrint.Condition = conditions
		}
		o := structs.Map(httpRulePrint)
		out = append(out, o)
	}
	return out
}
