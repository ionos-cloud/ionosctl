package lan

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

func K8sNodePoolLanAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nodepool",
		Resource:  "lan",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Kubernetes NodePool LAN",
		LongDesc: `Use this command to add a Node Pool LAN into an existing Node Pool.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id`,
		Example:    "ionosctl k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunK8sClusterNodePoolLanIds,
		CmdRun:     cloudapiv6cmds.RunK8sNodePoolLanAdd,
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
	cmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be attached to worker Nodes", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	cmd.AddStringSliceFlag(cloudapiv6.ArgNetwork, "", nil, "Slice of IPv4 or IPv6 CIDRs to be routed via the interface. Must contain same number of arguments as --gateway-ip flag")
	cmd.AddStringSliceFlag(cloudapiv6.ArgGatewayIp, "", nil, "Slice of IPv4 or IPv6 Gateway IPs for the routes. Must contain same number of arguments as --network flag")
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultK8sNodePoolLanCols, tabheaders.ColsMessage(defaultK8sNodePoolLanCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodePoolLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
