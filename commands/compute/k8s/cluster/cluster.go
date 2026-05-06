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
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Kubernetes Cluster Operations",
			Long:             "The sub-commands of `ionosctl compute k8s` allow you to perform Kubernetes Operations.",
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
