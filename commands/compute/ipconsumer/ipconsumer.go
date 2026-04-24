package ipconsumer

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allIpConsumerCols = []table.Column{
	{Name: "Ip", JSONPath: "ip", Default: true},
	{Name: "Mac", JSONPath: "mac"},
	{Name: "NicId", JSONPath: "nicId", Default: true},
	{Name: "ServerId", JSONPath: "serverId", Default: true},
	{Name: "ServerName", JSONPath: "serverName"},
	{Name: "DatacenterId", JSONPath: "datacenterId", Default: true},
	{Name: "DatacenterName", JSONPath: "datacenterName"},
	{Name: "K8sNodePoolId", JSONPath: "k8sNodePoolUuid", Default: true},
	{Name: "K8sClusterId", JSONPath: "k8sClusterUuid", Default: true},
}

func IpconsumerCmd() *core.Command {
	ipconsumerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipconsumer",
			Aliases:          []string{"ipc"},
			Short:            "Ip Consumer Operations",
			Long:             "The sub-command of `ionosctl compute ipconsumer` allows you to list information about where each IP address from an IpBlock is being used.",
			TraverseChildren: true,
		},
	}
	ipconsumerCmd.AddColsFlag(allIpConsumerCols)

	ipconsumerCmd.AddCommand(IpconsumerListCmd())

	return core.WithConfigOverride(ipconsumerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
