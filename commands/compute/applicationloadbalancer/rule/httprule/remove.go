package httprule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AlbRuleHttpRuleRemoveCmd() *core.Command {
	removeCmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "httprule",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove an HTTP Rule from an Application Load Balancer Forwarding Rule",
		LongDesc: `Use this command to remove an HTTP rule (identified by its --name) from a forwarding rule on an Application Load Balancer. Use --all to remove every HTTP rule from the listener at once.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Application Load Balancer Id
* Forwarding Rule Id
* Http Rule Name`,
		Example:    "ionosctl compute alb rule httprule remove --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME",
		PreCmdRun:  PreRunApplicationLoadBalancerRuleHttpRuleDelete,
		CmdRun:     RunAlbRuleHttpRuleRemove,
		InitClient: true,
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, "", "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AlbForwardingRulesIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgApplicationLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name of the HTTP rule to remove (as shown by `alb rule httprule list`).", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all HTTP Rules")

	return removeCmd
}
