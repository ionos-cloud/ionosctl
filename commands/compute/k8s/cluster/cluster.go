package cluster

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allK8sClusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "K8sVersion", JSONPath: "properties.k8sVersion", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "MaintenanceWindow", Default: true, Format: func(item map[string]any) any {
		mw, ok := table.Navigate(item, "properties.maintenanceWindow").(map[string]any)
		if !ok || mw == nil {
			return nil
		}
		return fmt.Sprintf("%s %s", mw["dayOfTheWeek"], mw["time"])
	}},
	{Name: "Public", JSONPath: "properties.public", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "NatGatewayIp", JSONPath: "properties.natGatewayIp"},
	{Name: "NodeSubnet", JSONPath: "properties.nodeSubnet"},
	{Name: "AvailableUpgradeVersions", JSONPath: "properties.availableUpgradeVersions"},
	{Name: "ViableNodePoolVersions", JSONPath: "properties.viableNodePoolVersions"},
	{Name: "S3Bucket", JSONPath: "properties.s3Buckets"},
	{Name: "ApiSubnetAllowList", JSONPath: "properties.apiSubnetAllowList"},
}

func K8sClusterCmd() *core.Command {
	k8sClusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "Kubernetes Cluster (control plane) Operations",
			Long: `Manage Managed Kubernetes clusters - the control plane layer.

A cluster is the IONOS-managed Kubernetes control plane (API server, scheduler,
etcd). It has no worker capacity of its own: worker Nodes live in node pools
(see ` + "`ionosctl compute k8s nodepool`" + `) that you attach afterwards.

Key properties of a cluster:
  * Kubernetes version    - the control-plane version; node pools must run this
                            version or lower (see --k8s-version).
  * maintenance window    - a weekly window (day + time) during which IONOS may
                            apply patches and minor upgrades.
  * public vs private     - a public cluster exposes the API server on a public
                            endpoint; a private cluster keeps the control plane
                            on your network and requires a NAT gateway IP and a
                            location.
  * API subnet allow list - CIDRs permitted to reach the Kubernetes API server.
  * S3 audit log bucket   - an Object Storage bucket that receives API-server
                            audit logs.

A cluster must be ACTIVE before you can create node pools in it, and must contain
no node pools before it can be deleted.`,
			TraverseChildren: true,
		},
	}
	k8sClusterCmd.AddColsFlag(allK8sClusterCols)

	k8sClusterCmd.AddCommand(K8sClusterListCmd())
	k8sClusterCmd.AddCommand(K8sClusterGetCmd())
	k8sClusterCmd.AddCommand(K8sClusterCreateCmd())
	k8sClusterCmd.AddCommand(K8sClusterUpdateCmd())
	k8sClusterCmd.AddCommand(K8sClusterDeleteCmd())

	return core.WithConfigOverride(k8sClusterCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
