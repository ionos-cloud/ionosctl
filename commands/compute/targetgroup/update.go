package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TargetGroupUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "targetgroup",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Target Group",
		LongDesc: `Update a Target Group's distribution, health-check, or HTTP health-check settings.

Only the flags you pass are changed; unspecified settings keep their current values. This command does NOT manage the targets (backend servers) themselves - use the ` + "`" + `target` + "`" + ` sub-commands (add/remove) for that.

Changes propagate to every ALB forwarding rule that references this group, so tightening the health check or switching --algorithm affects live traffic distribution.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Target Group Id`,
		Example: `# Rename a group
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --name new-name -w

# Switch the algorithm to source-IP stickiness
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --algorithm SOURCE_IP

# Retune the HTTP health check to match a body instead of a status code
ionosctl compute targetgroup update --targetgroup-id TARGET_GROUP_ID --match-type RESPONSE_BODY --response OK`,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Updated Target Group", "The name of the target group. Used only for display; does not need to be unique.")
	cmd.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "How traffic is distributed across targets. ROUND_ROBIN: served alternately, honoring weights. LEAST_CONNECTION: the target with the fewest active connections is served next. RANDOM: chosen by a consistent pseudo-random function. SOURCE_IP: the same client IP always reaches the same target (source-based stickiness).")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "The forwarding protocol. Only HTTP is currently supported by the API.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Connection Health Check] Maximum time in milliseconds to wait for a target to respond to a check. If a target also has --check-interval set, the smaller of the two values is used once the TCP connection is established.")
	cmd.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Connection Health Check] Interval in milliseconds between consecutive health checks. Default is 2000.")
	cmd.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Connection Health Check] Maximum number of reconnection attempts to a target after a connection failure before it is marked unhealthy. Valid range is 0 to 65535; default is 3.")
	cmd.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The request path (URL) the check sends to each target, e.g. /healthz. Default is '/.'.")
	cmd.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The HTTP method used for the health check request.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] What part of the target's reply decides health. STATUS_CODE: --response is matched against the HTTP status code. RESPONSE_BODY: --response is matched against the response body.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMatchType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The value a target must return to be considered healthy. Interpreted per --match-type: a status code (e.g. 200) for STATUS_CODE, or expected body text for RESPONSE_BODY.")
	cmd.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Treat --response as a regular expression when matching the response body, instead of a literal value. Default is false.")
	cmd.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Invert the match: the target is healthy when --response does NOT match. Default is false.")

	return cmd
}
