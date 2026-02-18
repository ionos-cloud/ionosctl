package networkloadbalancer

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer",
		LongDesc: `Use this command to update a specified Network Load Balancer from a Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id`,
		Example:    `ionosctl networkloadbalancer update --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID --name NAME`,
		PreCmdRun:  cloudapiv6cmds.PreRunDcNetworkLoadBalancerIds,
		CmdRun:     cloudapiv6cmds.RunNetworkLoadBalancerUpdate,
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
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of IP addresses of the Network Load Balancer")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Network Load Balancer update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer update [seconds]")

	return cmd
}
