package flowlog

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerFlowLogUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "flowlog",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an Application Load Balancer FlowLog",
		LongDesc: `Use this command to update a flow log on an Application Load Balancer. You can rename it, change which connections it captures (--action), the traffic direction (--direction), or the destination bucket (--s3bucket). Only the flags you provide are changed.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Application Load Balancer FlowLog Id`,
		Example: `# Change a flow log to capture only rejected connections
ionosctl compute applicationloadbalancer flowlog update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID --action REJECTED --wait`,
		PreCmdRun:  PreRunDcApplicationLoadBalancerFlowLogIds,
		CmdRun:     RunApplicationLoadBalancerFlowLogUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.Name(), cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgFlowLogId, cloudapiv6.ArgIdShort, "", cloudapiv6.FlowLogId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFlowLogId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancerFlowLogsIds(viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the Application Load Balancer FlowLog.")
	cmd.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "", "Which connections to log by their disposition: ACCEPTED (allowed), REJECTED (denied), or ALL.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "", "Which traffic direction to log relative to the ALB: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL (both).")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "The name of an existing IONOS Object Storage (S3) bucket that will receive the flow log records.")
	cmd.AddColsFlag(allFlowLogCols)

	return cmd
}
