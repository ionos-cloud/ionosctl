package networkloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Network Load Balancer",
		LongDesc: `Use this command to delete a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id`,
		Example:    `ionosctl compute networkloadbalancer delete --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID`,
		PreCmdRun:  PreRunDcNetworkLoadBalancerDelete,
		CmdRun:     RunNetworkLoadBalancerDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Network Load Balancer deletion to be executed")
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Network Load Balancers.")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer deletion [seconds]")

	return cmd
}
