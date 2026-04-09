package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
)

var clusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "DnsName", JSONPath: "metadata.dnsName", Default: true},
	{Name: "PostgresVersion", JSONPath: "properties.version", Default: true},
	{Name: "Instances", JSONPath: "properties.instances.count", Default: true},
	{Name: "Ram", JSONPath: "properties.instances.ram", Default: true},
	{Name: "Cores", JSONPath: "properties.instances.cores", Default: true},
	{Name: "StorageSize", JSONPath: "properties.instances.storageSize", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "SyncMode", JSONPath: "properties.replicationMode", Default: true},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "ConnectionPooler", JSONPath: "properties.connectionPooler"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "BackupLocation", JSONPath: "properties.backupLocation"},
	{Name: "LogsEnabled", JSONPath: "properties.logsEnabled"},
	{Name: "MetricsEnabled", JSONPath: "properties.metricsEnabled"},
	{Name: "DatacenterId", JSONPath: "properties.connection.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connection.lanId"},
	{Name: "Cidr", JSONPath: "properties.connection.primaryInstanceAddress"},
	{Name: "DbUsername", JSONPath: "properties.credentials.username"},
	{Name: "DbDatabase", JSONPath: "properties.credentials.database"},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres-v2 cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(clusterCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(clusterCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

// Clusters returns all clusters matching the given filters
func Clusters(fs ...Filter) (psqlv2.ClusterReadList, error) {
	req := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return psqlv2.ClusterReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return psqlv2.ClusterReadList{}, err
	}
	return ls, nil
}

type Filter func(request psqlv2.ApiClustersGetRequest) (psqlv2.ApiClustersGetRequest, error)
