package lan

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

func K8sNodePoolLanRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "nodepool",
		Resource:  "lan",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Kubernetes NodePool LAN",
		LongDesc: `This command removes a Kubernetes Node Pool LAN from a Node Pool.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* Lan Id`,
		Example:    "ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunK8sClusterNodePoolLanRemove,
		CmdRun:     cloudapiv6cmds.RunK8sNodePoolLanRemove,
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
	cmd.AddIntFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, 0, "The unique LAN Id of existing LANs to be detached from worker Nodes", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all FK8s Nodepool Lans.")

	return cmd
}
