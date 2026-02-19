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

func ApplicationLoadBalancerForwardingRuleUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to update a specified Application Load Balancer Forwarding Rule from a Application Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id`,
		Example:    "ionosctl compute alb rule update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME",
		PreCmdRun:  PreRunDcApplicationLoadBalancerForwardingRuleIds,
		CmdRun:     RunApplicationLoadBalancerForwardingRuleUpdate,
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
	cmd.AddStringFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.Name(), cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the Application Load Balancer forwarding rule.")
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Listening (inbound) IP.")
	cmd.AddIntFlag(cloudapiv6.ArgListenerPort, "", 8080, "Listening (inbound) port number; valid range is 1 to 65535.")
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 50, "The maximum time in milliseconds to wait for the client to acknowledge or send data; default is 50,000 (50 seconds).")
	cmd.AddStringSliceFlag(cloudapiv6.ArgServerCertificates, "", []string{""}, "Server Certificates")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, cloudapiv6.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.LbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultAlbForwardingRuleCols, tabheaders.ColsMessage(defaultAlbForwardingRuleCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
