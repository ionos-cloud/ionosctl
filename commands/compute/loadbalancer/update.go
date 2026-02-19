package loadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadBalancerUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "loadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Load Balancer",
		LongDesc: `Use this command to update the configuration of a specified Load Balancer.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Load Balancer Id`,
		Example:    "ionosctl compute loadbalancer update --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID",
		PreCmdRun:  PreRunDcLoadBalancerIds,
		CmdRun:     RunLoadBalancerUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.LoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Load Balancer")
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The IP of the Load Balancer")
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", cloudapiv6.DefaultDhcp, "Indicates if the Load Balancer will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for Load Balancer update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Load Balancer update [seconds]")

	return cmd
}
