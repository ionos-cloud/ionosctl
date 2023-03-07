package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultAlbRuleHttpRuleCols, printer.ColsMessage(allAlbRuleHttpRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(albRuleHttpRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = albRuleHttpRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbRuleHttpRuleCols, cobra.ShellCompDirectiveNoFileComp
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
		Example:    listApplicationLoadBalancerForwardingRuleHttpExample,
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunAlbRuleHttpRuleList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

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
		Example:    addApplicationLoadBalancerForwardingRuleHttpExample,
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRule,
		CmdRun:     RunAlbRuleHttpRuleAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The unique name of the Application Load Balancer HTTP rule.", core.RequiredFlagOption())
	add.AddStringFlag(cloudapiv6.ArgType, "", "", "Type of the HTTP rule.", core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, "", "", "he ID of the target group; mandatory and only valid for FORWARD actions.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddBoolFlag(cloudapiv6.ArgQuery, cloudapiv6.ArgQueryShort, false, "Default is false; valid only for REDIRECT actions.")
	add.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "www.ionos.com", "The location for redirecting; mandatory and valid only for REDIRECT actions.")
	add.AddIntFlag(cloudapiv6.ArgStatusCode, "", 301, "Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgStatusCode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"301", "302", "303", "307", "308", "200", "503", "599"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgMessage, cloudapiv6.ArgMessageShort, "Application Down", "The response message of the request; mandatory for STATIC actions.")
	add.AddStringFlag(cloudapiv6.ArgContentType, "", "application/json", "Valid only for STATIC actions.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgContentType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"application/json", "text/html"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgCondition, cloudapiv6.ArgConditionShort, "STARTS_WITH", "Matching rule for the HTTP rule condition attribute; mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCondition, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgConditionType, cloudapiv6.ArgConditionTypeShort, "HEADER", "Type of the HTTP rule condition.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgConditionType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "Specifies whether the condition is negated or not; the default is False.")
	add.AddStringFlag(cloudapiv6.ArgConditionKey, cloudapiv6.ArgConditionKeyShort, "forward-at", "Must be null when type is PATH, METHOD, HOST, or SOURCE_IP. Key can only be set when type is COOKIES, HEADER, or QUERY.")
	add.AddStringFlag(cloudapiv6.ArgConditionValue, cloudapiv6.ArgConditionValueShort, "Friday", "Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS; should be a valid CIDR if provided and if type is SOURCE_IP.")
	add.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule creation to be executed")
	add.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule creation [seconds]")
	add.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

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
* Http Rule Name`,
		Example:    removeApplicationLoadBalancerForwardingRuleHttpExample,
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRuleDelete,
		CmdRun:     RunAlbRuleHttpRuleRemove,
		InitClient: true,
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A name of that Application Load Balancer Http Rule", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all HTTP Rules")
	removeCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule deletion to be executed")
	removeCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule deletion [seconds]")
	removeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	return albRuleHttpRuleCmd
}

func PreRunApplicationLoadBalancerRuleHttpRule(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId,
		cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgName, cloudapiv6.ArgType)
}

func PreRunApplicationLoadBalancerRuleHttpRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgName},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgAll},
	)
}

func RunAlbRuleHttpRuleList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	c.Printer.Verbose("Getting HttpRules")
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
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
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	var httpRuleItems []ionoscloud.ApplicationLoadBalancerHttpRule
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	c.Printer.Verbose("Getting HttpRules from ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	ngOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
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
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
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
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
		resp, err := RemoveAllHTTPRules(c)
		if err != nil {
			return err
		}
		return c.Printer.Print(getAlbRuleHttpRulePrint(resp, c, nil))
	} else {
		c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Printer.Verbose("ApplicationLoadBalancer ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		c.Printer.Verbose("ForwardingRule ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "remove forwarding rule http rule"); err != nil {
			return err
		}
		c.Printer.Verbose("Getting HttpRules")
		frOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
			queryParams,
		)
		if err != nil {
			return err
		}
		c.Printer.Verbose("Removing the HTTP Rule from the existing HTTP Rules")
		proper, err := getRuleHttpRulesRemove(c, frOld)
		if err != nil {
			return err
		}
		c.Printer.Verbose("Updating ForwardingRule with the new HTTP Rules")
		_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
			proper,
			queryParams,
		)
		if resp != nil {
			c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getAlbRuleHttpRulePrint(resp, c, nil))
	}
}

func RemoveAllHTTPRules(c *core.CommandConfig) (*resources.Response, error) {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return nil, err
	}
	queryParams := listQueryParams.QueryParams
	_ = c.Printer.Warn("Forwarding Rule HTTP Rules to be deleted:")
	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		cloudapiv6.ParentResourceQueryParams,
	)
	if err != nil {
		return nil, err
	}
	if propertiesOk, ok := applicationLoadBalancerRules.GetPropertiesOk(); ok && propertiesOk != nil {
		if httpRulesOk, ok := propertiesOk.GetHttpRulesOk(); ok && httpRulesOk != nil {
			for _, httpRuleOk := range *httpRulesOk {
				if nameOk, ok := httpRuleOk.GetNameOk(); ok && nameOk != nil {
					_ = c.Printer.Warn("Forwarding Rule HTTP Rule Name: " + *nameOk)
				}
				if typeOk, ok := httpRuleOk.GetTypeOk(); ok && typeOk != nil {
					_ = c.Printer.Warn("Forwarding Rule HTTP Rule Type: " + *typeOk)
				}
			}

		}
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Forwarding Rule HTTP Rules"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Forwarding Rule HTTP Rules...")
		propertiesOk.SetHttpRules([]ionoscloud.ApplicationLoadBalancerHttpRule{})
		_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
			&resources.ApplicationLoadBalancerForwardingRuleProperties{
				ApplicationLoadBalancerForwardingRuleProperties: *propertiesOk,
			},
			queryParams,
		)
		if resp != nil {
			c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
			c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return nil, err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return nil, err
		}
	}
	return resp, err
}

func getRuleHttpRuleInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRule {
	// Set Application Load Balancer HTTP Rule Properties
	httprule := resources.ApplicationLoadBalancerHttpRule{}
	httprule.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	httprule.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	c.Printer.Verbose("Property Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)))
	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)), "FORWARD") {
		httprule.SetTargetGroup(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		c.Printer.Verbose("Property TargetGroup set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	}
	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)), "REDIRECT") {
		httprule.SetLocation(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
		c.Printer.Verbose("Property Location set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
		httprule.SetDropQuery(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgQuery)))
		c.Printer.Verbose("Property DropQuery set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgQuery)))
		httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
		c.Printer.Verbose("Property StatusCode set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
	}
	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType)), "STATIC") {
		httprule.SetResponseMessage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMessage)))
		c.Printer.Verbose("Property ResponseMessage set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMessage)))
		httprule.SetContentType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))
		c.Printer.Verbose("Property ContentType set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)) {
			httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
			c.Printer.Verbose("Property StatusCode set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
		} else {
			httprule.SetStatusCode(503)
			c.Printer.Verbose("Property StatusCode set with the default value: %v", 503)
		}
	}
	httpRuleCondition := getRuleHttpRuleConditionInfo(c)
	httprule.SetConditions([]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
		httpRuleCondition.ApplicationLoadBalancerHttpRuleCondition,
	})
	c.Printer.Verbose("Setting Condition to HttpRule")
	return httprule
}

func getRuleHttpRuleConditionInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRuleCondition {
	// Set Application Load Balancer HTTP Rule Condition Properties
	httpRuleCondition := resources.ApplicationLoadBalancerHttpRuleCondition{}
	httpRuleCondition.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)))
	c.Printer.Verbose("Property Condition Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)))
	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "SOURCE_IP") {
		httpRuleCondition.SetCondition(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
		c.Printer.Verbose("Property Condition set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
	}
	httpRuleCondition.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	c.Printer.Verbose("Property Condition Negate set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "COOKIES") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "HEADER") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "QUERY") {
		httpRuleCondition.SetKey(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
		c.Printer.Verbose("Property Condition Key set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
	}
	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)), "EXISTS") {
		httpRuleCondition.SetValue(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
		c.Printer.Verbose("Property Condition Value set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
	}
	return httpRuleCondition
}

func getRuleHttpRulesRemove(c *core.CommandConfig, frOld *resources.ApplicationLoadBalancerForwardingRule) (*resources.ApplicationLoadBalancerForwardingRuleProperties, error) {
	httpRuleItems := make([]ionoscloud.ApplicationLoadBalancerHttpRule, 0)
	if properties, ok := frOld.GetPropertiesOk(); ok && properties != nil {
		c.Printer.Verbose("Getting Properties from the Forwarding Rule")
		if httpRules, ok := properties.GetHttpRulesOk(); ok && httpRules != nil {
			c.Printer.Verbose("Getting HTTP Rules from the Forwarding Rule Properties")
			for _, httpRuleItem := range *httpRules {
				removeName := false
				if nameOk, ok := httpRuleItem.GetNameOk(); ok && nameOk != nil {
					if *nameOk == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
						removeName = true
						c.Printer.Verbose("Found HTTP Rule with Name: %v", *nameOk)
					}
				}
				// If the Http rule with the unique name is found, continue.
				// If not, add it to the Forwarding Rule properties.
				if removeName {
					continue
				} else {
					httpRuleItems = append(httpRuleItems, httpRuleItem)
				}
			}
		}
	}
	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &httpRuleItems,
		},
	}, nil
}

// Output Printing

var (
	defaultAlbRuleHttpRuleCols = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Condition"}
	allAlbRuleHttpRuleCols     = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Location", "StatusCode", "ResponseMessage", "ContentType", "Condition"}
)

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
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getAlbRuleHttpRulesKVMaps(ss)
			r.Columns = getAlbRuleHttpRulesCols(core.GetFlagName(c.Resource, constants.ArgCols), c.Printer.GetStderr())
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
