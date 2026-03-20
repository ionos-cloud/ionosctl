package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerForwardingRuleCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a Application Load Balancer Forwarding Rule in a specified Application Load Balancer. You can also set Health Check Settings for Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Listener Ip
* Listener Port`,
		Example:    "ionosctl compute applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --listener-ip LISTENER_IP --listener-port LISTENER_PORT",
		PreCmdRun:  PreRunApplicationLoadBalancerForwardingRuleCreate,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleCreate,
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
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name of the Application Load Balancer forwarding rule.")
	cmd.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Balancing protocol.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Listening (inbound) IP. It must be assigned to the listener NIC of Application Load Balancer.", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening (inbound) port number; valid range is 1 to 65535.", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 50, "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).")
	cmd.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule creation [seconds]")
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultAlbForwardingRuleCols, tabheaders.ColsMessage(defaultAlbForwardingRuleCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
