package targetgroup

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TargetGroupCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "targetgroup",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Target Group",
		LongDesc: `Use this command to create a Target Group.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.`,
		Example:    `ionosctl targetgroup create --name TARGET_GROUP_NAME`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunTargetGroupCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Target Group", "The name of the target group.")
	cmd.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Balancing algorithm.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Balancing protocol")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Health Check] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.")
	cmd.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] The interval in milliseconds between consecutive health checks; default is 2000.")
	cmd.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection attempts.")
	cmd.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The path (destination URL) for the HTTP health check request; the default is /.")
	cmd.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The method for the HTTP health check.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] Match Type for the HTTP health check.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMatchType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The response returned by the request, depending on the match type.")
	cmd.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Regex for the HTTP health check.")
	cmd.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Negate for the HTTP health check.")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group creation to be executed.")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group creation [seconds].")

	return cmd
}
