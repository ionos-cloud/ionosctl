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
		ShortDesc: "Add an HTTP Rule to an Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to add an HTTP rule to a forwarding rule (listener) on an Application Load Balancer. An HTTP rule matches incoming requests and then performs one action based on its --type:

  * FORWARD  - proxy matching requests to a backend target group (--targetgroup-id required).
  * REDIRECT - reply with an HTTP redirect to --location using --status-code (301/302/303/307/308); --query controls whether the original query string is dropped.
  * STATIC   - reply directly from the balancer with --status-code, --message and --content-type, without contacting any backend.

Matching is controlled by the condition flags: --condition-type selects which part of the request to inspect, --condition-key narrows it (e.g. a header name), --condition is the comparison operator, and --condition-value is what to compare against. Use --negate to invert the match. A rule with no conditions always matches and is useful as a default; rules are evaluated in order within the listener.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name
* Http Rule Type`,
		Example: `# FORWARD every request to a target group (no conditions = catch-all)
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "to-backend" --type FORWARD --targetgroup-id TARGETGROUP_ID

# FORWARD only requests whose path starts with /api
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "api-route" --type FORWARD --targetgroup-id TARGETGROUP_ID --condition-type PATH --condition STARTS_WITH --condition-value /api

# REDIRECT all traffic to an HTTPS host with a 301
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "force-https" --type REDIRECT --location https://www.ionos.com --status-code 301

# STATIC maintenance page returned directly by the balancer
ionosctl compute alb rule httprule add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --name "maintenance" --type STATIC --status-code 503 --content-type text/html --message "<h1>Under maintenance</h1>"`,
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
	cmd.AddStringFlag(cloudapiv6.ArgConditionType, cloudapiv6.ArgConditionTypeShort, "HEADER", "Selects which part of the incoming HTTP request the condition inspects. Possible values: HEADER, PATH, QUERY, METHOD, HOST, COOKIE, SOURCE_IP")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgConditionType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEADER", "PATH", "QUERY", "METHOD", "HOST", "COOKIE", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgConditionKey, cloudapiv6.ArgConditionKeyShort, "Accept", "Narrows the condition to a specific named entry within the selected element, e.g. the header name \"Accept\" when condition-type=HEADER. Only valid for HEADER, COOKIE and QUERY; must be empty for PATH, METHOD, HOST and SOURCE_IP.")
	cmd.AddStringFlag(cloudapiv6.ArgCondition, cloudapiv6.ArgConditionShort, "EQUALS", "The comparison operator applied between the selected request element and --condition-value. Possible values: EXISTS, CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH. Mandatory for HEADER, PATH, QUERY, METHOD, HOST and COOKIE types; must be empty when condition-type is SOURCE_IP.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCondition, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"EXISTS", "CONTAINS", "EQUALS", "MATCHES", "STARTS_WITH", "ENDS_WITH"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgConditionValue, cloudapiv6.ArgConditionValueShort, "application/json", "The value compared against the selected request element, e.g. \"application/json\" with condition=EQUALS, condition-type=HEADER, condition-key=Accept. Mandatory for CONTAINS, EQUALS, MATCHES, STARTS_WITH, ENDS_WITH, and for condition-type SOURCE_IP (where it must be a valid CIDR); must be empty when condition is EXISTS.")
	cmd.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "Inverts the condition so the rule matches when the condition does NOT hold. Default is false.")

	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The unique name of the Application Load Balancer HTTP rule (also used to reference it when removing).", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagType, "", "", "The action the rule performs on matching requests: FORWARD (proxy to a target group), STATIC (reply directly from the balancer), or REDIRECT (send an HTTP redirect).", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"FORWARD", "STATIC", "REDIRECT"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, "", "", "The target group whose backend servers matching requests are proxied to. Mandatory and only valid for --type FORWARD.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgQuery, cloudapiv6.ArgQueryShort, false, "When true, drops the query string from the redirect target so the redirect URI carries no query parameters. Default is false; valid only for --type REDIRECT.")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "www.ionos.com", "The URL/host clients are redirected to. Mandatory and only valid for --type REDIRECT.")
	cmd.AddIntFlag(cloudapiv6.ArgStatusCode, "", 301, "The HTTP status code returned to the client. Only valid for REDIRECT and STATIC actions. For REDIRECT: 301, 302, 303, 307 or 308 (default 301). For STATIC: any value in the range 200-599 (API default 503).")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgStatusCode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"301", "302", "303", "307", "308", "200", "503", "599"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgMessage, cloudapiv6.ArgMessageShort, "Application Down", "The response body the balancer returns. Mandatory and only used for --type STATIC.")
	cmd.AddStringFlag(cloudapiv6.ArgContentType, "", "application/json", "The Content-Type header of the static response, e.g. application/json or text/html. Only valid for --type STATIC.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgContentType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"application/json", "text/html"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.Flags().SortFlags = false

	return cmd
}
