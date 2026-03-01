package httprule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AlbRuleHttpRuleAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
		Example:    "ionosctl compute alb rule http add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME --type TYPE",
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRule,
		CmdRun:     RunAlbRuleHttpRuleAdd,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.Name(), cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})

	// see https://github.com/ionos-cloud/ionosctl/issues/263#issuecomment-1485258399
	cmd.AddStringFlag(cloudapiv6.ArgConditionType, cloudapiv6.ArgConditionTypeShort, "HEADER", "selects which element in the incoming HTTP request is used for the rule. Possible values HEADER, PATH, QUERY, METHOD, HOST, COOKIE, SOURCE _IP")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgConditionType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgConditionKey, cloudapiv6.ArgConditionKeyShort, "Accept", "selects which entry in the selected HTTP element is used for the rule. For example, \"Accept\" for condition-type=HEADER. Not applicable for HOST, PATH or SOURCE_IP")
	cmd.AddStringFlag(cloudapiv6.ArgCondition, cloudapiv6.ArgConditionShort, "EQUALS", "comparison rule for condition-value and the element selected with condition-type and condition-key. Possible values: EXISTS, CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH. mandatory for HEADER, PATH, QUERY, METHOD, HOST, and COOKIE types; must be null when type is SOURCE_IP.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCondition, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgConditionValue, cloudapiv6.ArgConditionValueShort, "application/json", "value compared with the selected HTTP element. For example \"application/json\" in combination with condition=EQUALS, condition-type=HEADER, condition-key=Accept would be valid. Not applicable for condition EXISTS. Mandatory for conditions CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH; must be null when condition is EXISTS")
	cmd.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "Specifies whether the condition is negated or not")

	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The unique name of the Application Load Balancer HTTP rule.", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagType, "", "", "Type of the HTTP rule.", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, "", "", "he ID of the target group; mandatory and only valid for FORWARD actions.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgQuery, cloudapiv6.ArgQueryShort, false, "Default is false; valid only for REDIRECT actions.")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "www.ionos.com", "The location for redirecting; mandatory and valid only for REDIRECT actions.")
	cmd.AddIntFlag(cloudapiv6.ArgStatusCode, "", 301, "Valid only for REDIRECT and STATIC actions. For REDIRECT actions, default is 301 and possible values are 301, 302, 303, 307, and 308. For STATIC actions, default is 503 and valid range is 200 to 599.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgStatusCode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"301", "302", "303", "307", "308", "200", "503", "599"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgMessage, cloudapiv6.ArgMessageShort, "Application Down", "The response message of the request; mandatory for STATIC actions.")
	cmd.AddStringFlag(cloudapiv6.ArgContentType, "", "application/json", "Valid only for STATIC actions.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgContentType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"application/json", "text/html"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Http Rule creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Http Rule creation [seconds]")

	cmd.Command.Flags().SortFlags = false

	return cmd
}
