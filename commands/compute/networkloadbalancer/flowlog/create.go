package flowlog

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerFlowLogCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer FlowLog",
		LongDesc: `Use this command to create a Network Load Balancer FlowLog in a specified Network Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Bucket Name`,
		Example:    `ionosctl networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --s3bucket BUCKET_NAME`,
		PreCmdRun:  cloudapiv6cmds.PreRunNetworkLoadBalancerFlowLogCreate,
		CmdRun:     cloudapiv6cmds.RunNetworkLoadBalancerFlowLogCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed FlowLog", "The name for the FlowLog")
	cmd.AddStringFlag(cloudapiv6.ArgAction, cloudapiv6.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAction, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgDirection, cloudapiv6.ArgDirectionShort, "BIDIRECTIONAL", "Specifies the traffic Direction pattern")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDirection, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3Bucket, cloudapiv6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Network Load Balancer FlowLog creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer FlowLog creation [seconds]")
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
