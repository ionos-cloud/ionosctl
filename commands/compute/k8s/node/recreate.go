package node

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

func K8sNodeRecreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "recreate",
		Aliases:   []string{"r"},
		ShortDesc: "Recreate a Kubernetes Node",
		LongDesc: `You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id`,
		Example:    "ionosctl k8s node recreate --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunK8sClusterNodesIds,
		CmdRun:     cloudapiv6cmds.RunK8sNodeRecreate,
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
	cmd.AddUUIDFlag(cloudapiv6.ArgK8sNodeId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sNodeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId)),
			viper.GetString(core.GetFlagName(cmd.NS, constants.FlagNodepoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
