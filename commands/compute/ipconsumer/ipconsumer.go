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
			Use:     "ipconsumer",
			Aliases: []string{"ipc"},
			Short:   "See which resources are using the IPs of a reserved IP block",
			Long: `The ` + "`ionosctl compute ipconsumer`" + ` command answers the question "what is currently using the IP addresses I reserved?".

An IP block (see ` + "`ionosctl compute ipblock`" + `) is a set of reserved public IPv4 addresses. Each address can be assigned to a NIC, which lives on a server inside a datacenter, or be used by a load balancer / Kubernetes node pool. This command lists, for every address in a given block, the consuming resource: the NIC, its MAC, the owning server and datacenter, and any Kubernetes cluster / node pool.

Use it to audit reservations, find which server is holding an address you want to free, or confirm a block is unused before deleting it. It complements ` + "`ionosctl compute ipblock`" + `, which manages the reservations themselves. This command is read-only.`,
			TraverseChildren: true,
		},
	}
	ipconsumerCmd.AddColsFlag(allIpConsumerCols)

	ipconsumerCmd.AddCommand(IpconsumerListCmd())

	return core.WithConfigOverride(ipconsumerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
