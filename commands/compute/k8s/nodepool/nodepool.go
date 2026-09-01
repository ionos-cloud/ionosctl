package nodepool

import (
	"fmt"

	nplan "github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/nodepool/lan"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allK8sNodePoolCols = []table.Column{
	{Name: "NodePoolId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "K8sVersion", JSONPath: "properties.k8sVersion", Default: true},
	{Name: "NodeCount", JSONPath: "properties.nodeCount", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.datacenterId", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "CpuFamily", JSONPath: "properties.cpuFamily"},
	{Name: "ServerType", JSONPath: "properties.serverType"},
	{Name: "StorageType", JSONPath: "properties.storageType"},
	{Name: "LanIds", JSONPath: "properties.lans.*.id"},
	{Name: "CoresCount", JSONPath: "properties.coresCount"},
	{Name: "RamSize", JSONPath: "properties.ramSize"},
	{Name: "AvailabilityZone", JSONPath: "properties.availabilityZone"},
	{Name: "StorageSize", JSONPath: "properties.storageSize"},
	{Name: "MaintenanceWindow", Format: func(item map[string]any) any {
		mw, ok := table.Navigate(item, "properties.maintenanceWindow").(map[string]any)
		if !ok || mw == nil {
			return nil
		}
		return fmt.Sprintf("%s %s", mw["dayOfTheWeek"], mw["time"])
	}},
	{Name: "AutoScaling", JSONPath: "properties.autoScaling"},
	{Name: "PublicIps", JSONPath: "properties.publicIps"},
	{Name: "AvailableUpgradeVersions", JSONPath: "properties.availableUpgradeVersions"},
	{Name: "Annotations", JSONPath: "properties.annotations"},
	{Name: "Labels", JSONPath: "properties.labels"},
	{Name: "ClusterId", JSONPath: "href"},
}

func K8sNodePoolCmd() *core.Command {
	k8sNodePoolCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "nodepool",
			Aliases: []string{"np"},
			Short:   "Kubernetes NodePool (worker Nodes) Operations",
			Long: `Manage Managed Kubernetes node pools - the worker-Node layer.

A node pool is a group of worker Nodes with identical hardware, all provisioned
into a single Data Center (--datacenter-id). Every Node in the pool shares the
same shape: --cores, --ram, --storage-size, --storage-type, --cpu-family and
Kubernetes version. A cluster may have many node pools (e.g. one per Data Center,
or pools of different sizes for different workloads).

Sizing and scaling:
  * --node-count sets a fixed number of Nodes.
  * Autoscaling (min/max node count) lets Managed Kubernetes grow and shrink the
    pool automatically between a lower and upper bound; it is configured via
    ` + "`nodepool update`" + `.

Version skew: a node pool's Kubernetes version must be less than or equal to its
cluster's version, and can only be upgraded (never downgraded). Upgrade the
cluster first, then the pool.

Networking: worker Nodes can attach to existing LANs in the Data Center (see
--lan-ids and the ` + "`nodepool lan`" + ` sub-commands for per-LAN routes).

The parent cluster must be ACTIVE before a node pool can be created.`,
			TraverseChildren: true,
		},
	}
	k8sNodePoolCmd.AddColsFlag(allK8sNodePoolCols)

	k8sNodePoolCmd.AddCommand(K8sNodePoolListCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolGetCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolCreateCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolUpdateCmd())
	k8sNodePoolCmd.AddCommand(K8sNodePoolDeleteCmd())
	k8sNodePoolCmd.AddCommand(nplan.K8sNodePoolLanCmd())

	return core.WithConfigOverride(k8sNodePoolCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
