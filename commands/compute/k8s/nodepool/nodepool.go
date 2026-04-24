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
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Kubernetes NodePool Operations",
			Long:             "The sub-commands of `ionosctl compute k8s nodepool` allow you to list, get, create, update, delete Kubernetes NodePools.",
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
