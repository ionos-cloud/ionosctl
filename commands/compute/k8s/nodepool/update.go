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
		ShortDesc: "Update a Kubernetes NodePool (worker Nodes)",
		LongDesc: `Update a node pool: resize it, enable or disable autoscaling, upgrade its
Kubernetes version, change labels/annotations, adjust the maintenance window,
attach LANs, or add reserved public IPs.

Autoscaling: set --min-node-count and --max-node-count together to let Managed
Kubernetes automatically scale the pool between those bounds. Set both to 0 to
disable autoscaling and return to a fixed --node-count.

Kubernetes version: --k8s-version is upgrade-only (downgrades are rejected) and
must stay <= the cluster's version. Upgrade the cluster first if needed.

Maintenance window: --maintenance-day and --maintenance-time define the weekly
window during which IONOS may recycle Nodes for patches/upgrades.

Reserved public IPs (--public-ips): IPs must be reserved in the same location as
the pool's Data Center. Provide one more IP than the pool's maximum node count -
the extra IP is used while Nodes are being rebuilt/rolled.

LANs: --lan-ids adds LANs to the ones already attached (it does not replace
them). To also set a route network and gateway IP on a LAN, use
` + "`ionosctl compute k8s nodepool lan add`" + ` per LAN.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the node pool returns to the AVAILABLE state.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id`,
		Example: `# Resize a node pool to 5 Nodes
ionosctl compute k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count 5

# Enable autoscaling between 2 and 8 Nodes and move the maintenance window to Sunday 03:00
ionosctl compute k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID \
  --min-node-count 2 --max-node-count 8 --maintenance-day Sunday --maintenance-time 03:00:00`,
		PreCmdRun:  PreRunK8sClusterNodePoolIds,
		CmdRun:     RunK8sNodePoolUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgK8sVersion, "", "", "New Kubernetes version for the worker Nodes, e.g. 1.29.5. Upgrade-only (downgrades are rejected) and must be <= the cluster's version")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sVersion,
		func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			clusterId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))
			nodepoolId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagNodepoolId))
			return completer.K8sNodePoolUpgradeVersions(clusterId, nodepoolId), cobra.ShellCompDirectiveNoFileComp
		})
	cmd.AddIntFlag(constants.FlagNodeCount, "", 1, "Fixed number of worker Nodes in the pool. Ignored while autoscaling is enabled (min/max node count)")
	cmd.AddIntFlag(cloudapiv6.ArgK8sMinNodeCount, "", 1, "Autoscaling lower bound: minimum number of Nodes the pool may scale in to. Set together with --max-node-count. Set both to 0 to disable autoscaling")
	cmd.AddIntFlag(cloudapiv6.ArgK8sMaxNodeCount, "", 1, "Autoscaling upper bound: maximum number of Nodes the pool may scale out to. Set together with --min-node-count. Set both to 0 to disable autoscaling")
	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Kubernetes labels for the pool's Nodes. Overwrites any existing labels. Format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Kubernetes annotations for the pool's Nodes. Overwrites any existing annotations. Format: --annotations KEY=VALUE,KEY=VALUE")
	cmd.AddStringFlag(cloudapiv6.ArgLabelKey, "", "", "Label key. Must be set together with --label-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgLabelValue, "", "", "Label value. Must be set together with --label-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sAnnotationKey, "", "", "Annotation key. Must be set together with --annotation-value", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sAnnotationValue, "", "", "Annotation value. Must be set together with --annotation-key", core.DeprecatedFlagOption("Use --labels, --annotations options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceDay, "", "", "Day of the week for the maintenance window, in English (e.g. Monday, Saturday). Set together with --maintenance-time")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgK8sMaintenanceDay, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagServerType, "", "", []string{"DedicatedCore", "VCPU"},
		"Compute-engine server type for the Nodes: 'DedicatedCore' (dedicated physical CPU cores) "+
			"or 'VCPU' (shared vCPU cores, typically cheaper)")
	cmd.AddStringFlag(cloudapiv6.ArgK8sMaintenanceTime, "", "", "Start time of the maintenance window in HH:mm:ss (24-hour) format, e.g. 08:00:00. Set together with --maintenance-day")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPublicIps, "", []string{}, "Reserved public IPs for the Nodes. IPs must be reserved in the same location as the pool's Data Center. Provide one more IP than the pool's maximum node count (the extra is used while rebuilding Nodes). Usage: --public-ips IP1,IP2")
	cmd.AddIntSliceFlag(cloudapiv6.ArgLanIds, "", []int{}, "IDs of existing LANs to attach to the worker Nodes. These are added to the LANs already attached, not replacing them. Usage: --lan-ids 1,2")
	cmd.AddBoolFlag(cloudapiv6.ArgDhcp, "", true, "Whether Nodes obtain an IP via DHCP on the LANs being attached. e.g. --dhcp=true, --dhcp=false")
	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sNodePoolId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
