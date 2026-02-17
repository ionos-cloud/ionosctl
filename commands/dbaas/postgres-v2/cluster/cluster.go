package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
	"github.com/spf13/cobra"
)

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "DnsName", "PostgresVersion", "Instances", "Ram", "Cores", "StorageSize", "State", "SyncMode"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "DnsName", "PostgresVersion", "Instances", "Ram", "Cores", "StorageSize", "State", "SyncMode",
		"MaintenanceDay", "MaintenanceTime", "BackupLocation", "DatacenterId", "LanId", "Cidr"}
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

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
