package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodePoolLanAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nodepool",
		Resource:  "lan",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Kubernetes NodePool LAN",
		LongDesc: `Attach an existing LAN to a node pool so its worker Nodes gain an interface on
that LAN. The LAN must already exist in the same Data Center as the pool.

Optionally add routes: --network gives destination CIDRs and --gateway-ip the
gateway each is reached through, so Nodes can route to networks behind that
gateway. The two flags are positional - the Nth --network is paired with the Nth
--gateway-ip, so they must have the same number of entries.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the node pool reaches the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id`,
		Example: `# Attach LAN 2 to a node pool with DHCP
ionosctl compute k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id 2

# Attach a LAN with a static route (10.0.0.0/24 via 10.1.5.16)
ionosctl compute k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id 2 --network 10.0.0.0/24 --gateway-ip 10.1.5.16`,
		PreCmdRun:  PreRunK8sClusterNodePoolLanIds,
		CmdRun:     RunK8sNodePoolLanAdd,
		InitClient: true,
	})
	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, "", "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "ID of an existing LAN (in the pool's Data Center) to attach to the worker Nodes", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Whether Nodes obtain an IP on this LAN via DHCP. e.g. --dhcp=true, --dhcp=false")
	cmd.AddStringSliceFlag(cloudapiv6.ArgNetwork, "", nil, "Destination IPv4/IPv6 CIDRs to route via this LAN. Paired positionally with --gateway-ip, so it must have the same number of entries")
	cmd.AddStringSliceFlag(cloudapiv6.ArgGatewayIp, "", nil, "Gateway IPs (IPv4/IPv6) for the corresponding --network routes. Paired positionally with --network, so it must have the same number of entries")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, table.ColsMessage(allK8sNodePoolLanCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allK8sNodePoolLanCols), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
