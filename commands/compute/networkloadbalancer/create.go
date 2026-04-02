package networkloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func NetworkLoadBalancerCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer",
		LongDesc: `Use this command to create a Network Load Balancer in a specified Virtual Data Center.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    `ionosctl compute networkloadbalancer create --datacenter-id DATACENTER_ID`,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunNetworkLoadBalancerCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Id of the listening LAN")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Id of the balanced private target LAN")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "Collection of IP addresses of the Network Load Balancer")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Collection of private IP addresses with subnet mask of the Network Load Balancer")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Network Load Balancer creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Network Load Balancer creation [seconds]")

	return cmd
}
