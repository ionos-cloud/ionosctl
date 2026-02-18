package nodepool

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodePoolUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "nodepool",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Kubernetes NodePool",
		LongDesc: `Use this command to update the number of worker Nodes, the minimum and maximum number of worker Nodes, the add labels, annotations, to update the maintenance day and time, to attach private LANs to a Node Pool within an existing Kubernetes Cluster. You can also add reserved public IP addresses to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. The array must contain one extra IP than maximum number of Nodes could be. The extra provided IP Will be used during rebuilding of Nodes.

You can wait for the Node Pool to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Note: If you want to attach multiple LANs to Node Pool, use ` + "`" + `--lan-ids=LAN_ID1,LAN_ID2` + "`" + ` flag. If you want to also set a Route Network, Route GatewayIp for LAN use ` + "`" + `ionosctl k8s nodepool lan add` + "`" + ` command for each LAN you want to add.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`,
		Example:    "ionosctl k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID",
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "The K8s version for the NodePool. K8s version downgrade is not supported")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sVersion,
		func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			clusterId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))
			nodepoolId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagNodepoolId))
			return completer.K8sNodePoolUpgradeVersions(clusterId, nodepoolId), cobra.ShellCompDirectiveNoFileComp
		})
	cmd.AddIntFlag(constants.FlagNodeCount, "", 1, "The number of worker Nodes that the NodePool should contain")
	cmd.AddIntFlag(cloudapiv6.ArgK8sMinNodeCount, "", 1, "The minimum number of worker Nodes that the managed NodePool can scale in. Should be set together with --max-node-count. Set to 0 to disable autoscaling.")
	cmd.AddIntFlag(cloudapiv6.ArgK8sMaxNodeCount, "", 1, "The maximum number of worker Nodes that the managed NodePool can scale out. Should be set together with --min-node-count. Set to 0 to disable autoscaling")
	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")
	cmd.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "Label key. Must be set together with --label-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", "Label value. Must be set together with --label-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sMaintenanceDay, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagServerType, "", "", []string{"DedicatedCore", "VCPU"},
		"The type of server for the Kubernetes node pool can be either"+
			"'DedicatedCore' (nodes with dedicated CPU cores) or 'VCPU' (nodes with shared CPU cores)."+
			"This selection corresponds to the server type for the compute engine.")
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPublicIps, "", []string{}, "Reserved public IP address to be used by the Nodes. IPs must be from same location as the Data Center used for the Node Pool. Usage: --public-ips IP1,IP2")
	cmd.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "Collection of LAN Ids of existing LANs to be attached to worker Nodes. It will be added to the existing LANs attached")
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Indicates if the Kubernetes Node Pool LANs will reserve an IP using DHCP. E.g.: --dhcp=true, --dhcp=false")
	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for the new NodePool to be in ACTIVE state")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.K8sTimeoutSeconds, "Timeout option for waiting for NodePool to be in ACTIVE state [seconds]")

	return cmd
}
