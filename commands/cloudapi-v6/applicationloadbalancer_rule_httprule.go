package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultAlbRuleHttpRuleCols = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Condition"}
	allAlbRuleHttpRuleCols     = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Location", "StatusCode", "ResponseMessage", "ContentType", "Condition"}
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
	globalFlags.StringSliceP(constants.FlagCols, "", defaultAlbRuleHttpRuleCols, tabheaders.ColsMessage(allAlbRuleHttpRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(albRuleHttpRuleCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = albRuleHttpRuleCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.FlagRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.FlagDepthDescription)

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
	add.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.FlagRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.FlagDataCenterId)), viper.GetString(core.GetFlagName(add.NS, cloudapiv6.FlagApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})

	// see https://github.com/ionos-cloud/ionosctl/issues/263#issuecomment-1485258399
	add.AddStringFlag(cloudapiv6.FlagConditionType, cloudapiv6.FlagConditionTypeShort, "HEADER", "selects which element in the incoming HTTP request is used for the rule. Possible values HEADER, PATH, QUERY, METHOD, HOST, COOKIE, SOURCE _IP")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagConditionType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.FlagConditionKey, cloudapiv6.FlagConditionKeyShort, "Accept", "selects which entry in the selected HTTP element is used for the rule. For example, \"Accept\" for condition-type=HEADER. Not applicable for HOST, PATH or SOURCE_IP")
	add.AddStringFlag(cloudapiv6.FlagCondition, cloudapiv6.FlagConditionShort, "EQUALS", "comparison rule for condition-value and the element selected with condition-type and condition-key. Possible values: EXISTS, CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH. mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagCondition, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.FlagConditionValue, cloudapiv6.FlagConditionValueShort, "application/json", "value compared with the selected HTTP element. For example \"application/json\" in combination with condition=EQUALS, condition-type=HEADER, condition-key=Accept would be valid. Not applicable for condition EXISTS. Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS")
	add.AddBoolFlag(cloudapiv6.FlagNegate, "", false, "Specifies whether the condition is negated or not")

	add.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "The unique name of the Application Load Balancer HTTP rule.", core.RequiredFlagOption())
	add.AddStringFlag(constants.FlagType, "", "", "Type of the HTTP rule.", core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.FlagTargetGroupId, "", "", "he ID of the target group; mandatory and only valid for FORWARD actions.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddBoolFlag(cloudapiv6.FlagQuery, cloudapiv6.FlagQueryShort, false, "Default is false; valid only for REDIRECT actions.")
	add.AddStringFlag(cloudapiv6.FlagLocation, cloudapiv6.FlagLocationShort, "www.ionos.com", "The location for redirecting; mandatory and valid only for REDIRECT actions.")
	add.AddIntFlag(cloudapiv6.FlagStatusCode, "", 301, "Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagStatusCode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"301", "302", "303", "307", "308", "200", "503", "599"}, cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.FlagMessage, cloudapiv6.FlagMessageShort, "Application Down", "The response message of the request; mandatory for STATIC actions.")
	add.AddStringFlag(cloudapiv6.FlagContentType, "", "application/json", "Valid only for STATIC actions.")
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagContentType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"application/json", "text/html"}, cobra.ShellCompDirectiveNoFileComp
	})

	add.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule creation to be executed")
	add.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule creation [seconds]")
	add.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.FlagDepthDescription)

	add.Command.Flags().SortFlags = false

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
	removeCmd.AddUUIDFlag(cloudapiv6.FlagDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.FlagApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.FlagDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.FlagRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.FlagRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.FlagDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.FlagApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.FlagName, cloudapiv6.FlagNameShort, "", "A name of that Application Load Balancer Http Rule", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(cloudapiv6.FlagAll, cloudapiv6.FlagAllShort, false, "Remove all HTTP Rules")
	removeCmd.AddBoolFlag(constants.FlagWaitForRequest, constants.FlagWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule deletion to be executed")
	removeCmd.AddIntFlag(constants.FlagTimeout, constants.FlagTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule deletion [seconds]")
	removeCmd.AddInt32Flag(cloudapiv6.FlagDepth, cloudapiv6.FlagDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.FlagDepthDescription)

	return core.WithConfigOverride(albRuleHttpRuleCmd, "compute", "")
}

func PreRunApplicationLoadBalancerRuleHttpRule(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.FlagDataCenterId,
		cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagRuleId, cloudapiv6.FlagName, constants.FlagType)
}

func PreRunApplicationLoadBalancerRuleHttpRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagRuleId, cloudapiv6.FlagName},
		[]string{cloudapiv6.FlagDataCenterId, cloudapiv6.FlagApplicationLoadBalancerId, cloudapiv6.FlagRuleId, cloudapiv6.FlagAll},
	)
}

func RunAlbRuleHttpRuleList(c *core.CommandConfig) error {
	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting HttpRules"))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := ng.GetPropertiesOk()
	if !ok || properties == nil {
		return errors.New("error getting rule properties")
	}

	httpRules, ok := properties.GetHttpRulesOk()
	if !ok || httpRules == nil {
		return errors.New("error getting rule http rules")
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApplicationLoadBalancerHTTPRule, *httpRules,
		tabheaders.GetHeaders(allAlbRuleHttpRuleCols, defaultAlbRuleHttpRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunAlbRuleHttpRuleAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	var httpRuleItems []ionoscloud.ApplicationLoadBalancerHttpRule

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting HttpRules from ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))

	ngOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Adding the new HttpRule to the existing HttpRules"))

	httpRuleItems = append(httpRuleItems, httpRuleNew.ApplicationLoadBalancerHttpRule)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating ForwardingRule with the new HttpRules"))

	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		&resources.ApplicationLoadBalancerForwardingRuleProperties{
			ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
				HttpRules: &httpRuleItems,
			},
		},
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApplicationLoadBalancerHTTPRule, httpRuleNew.ApplicationLoadBalancerHttpRule,
		tabheaders.GetHeaders(allAlbRuleHttpRuleCols, defaultAlbRuleHttpRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunAlbRuleHttpRuleRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))

		resp, err := RemoveAllHTTPRules(c)
		if err != nil {
			return err
		}
		if resp != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}

		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove forwarding rule http rule", viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting HttpRules"))

	frOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		queryParams,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing the HTTP Rule from the existing HTTP Rules"))

	proper, err := getRuleHttpRulesRemove(c, frOld)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating ForwardingRule with the new HTTP Rules"))

	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		proper,
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Application Load Balancer HTTP Rule successfully deleted"))

	return nil

}

func RemoveAllHTTPRules(c *core.CommandConfig) (*resources.Response, error) {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return nil, err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Forwarding Rule HTTP Rules to be deleted:"))

	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		cloudapiv6.ParentResourceQueryParams,
	)
	if err != nil {
		return nil, err
	}

	propertiesOk, ok := applicationLoadBalancerRules.GetPropertiesOk()
	if !ok || propertiesOk == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer Forwarding Rule properties")
	}

	httpRulesOk, ok := propertiesOk.GetHttpRulesOk()
	if !ok || httpRulesOk == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer HTTP Rules")
	}

	if len(*httpRulesOk) <= 0 {
		return nil, fmt.Errorf("no Application Load Balancer HTTP Rules found")
	}

	for _, httpRuleOk := range *httpRulesOk {
		if nameOk, ok := httpRuleOk.GetNameOk(); ok && nameOk != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Forwarding Rule HTTP Rule Name: %v", *nameOk))
		}

		if typeOk, ok := httpRuleOk.GetTypeOk(); ok && typeOk != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Forwarding Rule HTTP Rule Type: %v", *typeOk))
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Forwarding Rule HTTP Rules", viper.GetBool(constants.FlagForce)) {
		return nil, fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Forwarding Rule HTTP Rules..."))

	propertiesOk.SetHttpRules([]ionoscloud.ApplicationLoadBalancerHttpRule{})
	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagRuleId)),
		&resources.ApplicationLoadBalancerForwardingRuleProperties{
			ApplicationLoadBalancerForwardingRuleProperties: *propertiesOk,
		},
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Request Id: %v", request.GetId(resp)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return nil, err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return nil, err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Application Load Balancer HTTP Rules successfully deleted"))
	return resp, err
}

func getRuleHttpRuleInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRule {
	// Set Application Load Balancer HTTP Rule Properties
	httprule := resources.ApplicationLoadBalancerHttpRule{}

	httprule.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName))))

	httprule.SetType(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Type set: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagType))))

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "FORWARD") {
		httprule.SetTargetGroup(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagTargetGroupId)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property TargetGroup set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagTargetGroupId))))
	}

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "REDIRECT") {
		httprule.SetLocation(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLocation)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Location set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagLocation))))

		httprule.SetDropQuery(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagQuery)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property DropQuery set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagQuery))))

		httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagStatusCode)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property StatusCode set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagStatusCode))))
	}

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "STATIC") {
		httprule.SetResponseMessage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagMessage)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property ResponseMessage set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagMessage))))

		httprule.SetContentType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagContentType)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property ContentType set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagContentType))))

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagStatusCode)) {
			httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagStatusCode)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				"Property StatusCode set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.FlagStatusCode))))
		} else {
			httprule.SetStatusCode(503)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				"Property StatusCode set with the default value: %v", 503))
		}
	}

	httpRuleCondition := getRuleHttpRuleConditionInfo(c)
	httprule.SetConditions([]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
		httpRuleCondition.ApplicationLoadBalancerHttpRuleCondition,
	})

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Setting Condition to HttpRule"))

	return httprule
}

func getRuleHttpRuleConditionInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRuleCondition {
	// Set Application Load Balancer HTTP Rule Condition Properties
	httpRuleCondition := resources.ApplicationLoadBalancerHttpRuleCondition{}
	httpRuleCondition.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType)))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Condition Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType))))

	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType)), "SOURCE_IP") {
		httpRuleCondition.SetCondition(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagCondition)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Condition set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagCondition))))
	}

	httpRuleCondition.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagNegate)))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Condition Negate set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagNegate))))

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType)), "COOKIES") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType)), "HEADER") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionType)), "QUERY") {
		httpRuleCondition.SetKey(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionKey)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Condition Key set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionKey))))
	}

	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagCondition)), "EXISTS") {
		httpRuleCondition.SetValue(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionValue)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Property Condition Value set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagConditionValue))))
	}

	return httpRuleCondition
}

func getRuleHttpRulesRemove(c *core.CommandConfig, frOld *resources.ApplicationLoadBalancerForwardingRule) (*resources.ApplicationLoadBalancerForwardingRuleProperties, error) {
	httpRuleItems := make([]ionoscloud.ApplicationLoadBalancerHttpRule, 0)

	properties, ok := frOld.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer Forwarding Rule properties")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Properties from the Forwarding Rule"))

	httpRules, ok := properties.GetHttpRulesOk()
	if !ok || httpRules == nil {
		return nil, fmt.Errorf("coudl not get Application Load Balancer HTTP Rules")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting HTTP Rules from the Forwarding Rule Properties"))

	for _, httpRuleItem := range *httpRules {
		removeName := false

		if nameOk, ok := httpRuleItem.GetNameOk(); ok && nameOk != nil {
			if *nameOk == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagName)) {
				removeName = true
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Found HTTP Rule with Name: %v", *nameOk))
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

	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &httpRuleItems,
		},
	}, nil
}
