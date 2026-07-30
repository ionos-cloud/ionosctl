package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/spf13/cobra"
)

var clusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "DnsName", JSONPath: "metadata.dnsName", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Instances", JSONPath: "properties.instances.count"},
	{Name: "Cores", JSONPath: "properties.instances.cores"},
	{Name: "Ram", JSONPath: "properties.instances.ram"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "EvictionPolicy", JSONPath: "properties.evictionPolicy", Default: true},
	{Name: "PersistenceMode", JSONPath: "properties.persistenceMode"},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "SnapshotLocation", JSONPath: "properties.snapshot.location"},
	{Name: "RetentionDays", JSONPath: "properties.snapshot.retentionDays"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "LogsEnabled", JSONPath: "properties.logsEnabled"},
	{Name: "MetricsEnabled", JSONPath: "properties.metricsEnabled"},
	{Name: "DatacenterId", JSONPath: "properties.connection.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connection.lanId"},
	{Name: "Cidr", JSONPath: "properties.connection.primaryInstanceAddress"},
	{Name: "Username", JSONPath: "properties.credentials.username"},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "In-Memory DB Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas in-memory-db-v2 cluster` allow you to manage the In-Memory DB Clusters under your account.",
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
func Clusters(fs ...Filter) (inmemorydb.ClusterReadList, error) {
	req := client.Must().InMemoryDBClientV2.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return inmemorydb.ClusterReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return inmemorydb.ClusterReadList{}, err
	}
	return ls, nil
}

type Filter func(request inmemorydb.ApiClustersGetRequest) (inmemorydb.ApiClustersGetRequest, error)
