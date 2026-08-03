package nic

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadBalancerNicDetachCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "loadbalancer",
		Resource:  "nic",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a NIC from a Load Balancer",
		LongDesc: `Use this command to remove a NIC from a Load Balancer's backend pool. The NIC stops receiving balanced traffic and no longer shares the balancer's public IP; the NIC and its server are otherwise unaffected.

Use --all to detach every NIC currently attached to the Load Balancer in one call.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can skip the confirmation prompt with ` + "`" + `--force` + "`" + `.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id`,
		Example: `# Detach a single NIC from a Load Balancer
ionosctl compute loadbalancer nic detach --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID

# Detach every NIC attached to the Load Balancer, skipping confirmation
ionosctl compute loadbalancer nic detach --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --all --force`,
		PreCmdRun:  PreRunNicDetach,
		CmdRun:     RunLoadBalancerNicDetach,
		InitClient: true,
	})
	cmd.AddColsFlag(allNicCols)
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNicId, cloudapiv6.ArgIdShort, "", cloudapiv6.NicId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNicId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedNicsIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgLoadBalancerId, "", "", cloudapiv6.LoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LoadbalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach every NIC currently attached to the Load Balancer. When set, --nic-id is not required")

	return cmd
}
