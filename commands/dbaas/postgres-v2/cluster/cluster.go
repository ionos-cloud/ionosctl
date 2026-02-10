package cluster

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "Location", "PostgresVersion", "Instances", "Ram", "Cores", "StorageSize", "StorageType", "State", "SyncMode"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "Location", "PostgresVersion", "Instances", "Ram", "Cores", "StorageSize", "StorageType", "State", "SyncMode",
		"MaintenanceWindow", "DatacenterId", "LanId", "Cidr", "DnsName"}
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
