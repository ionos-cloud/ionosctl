package cluster

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

// validateBackupRetentionDays enforces the API's accepted retention range so a
// bad value fails fast with a clear message instead of an opaque server error.
func validateBackupRetentionDays(days int32) error {
	if days < 1 || days > 365 {
		return fmt.Errorf("--%s must be between 1 and 365 (got %d)", constants.FlagBackupRetentionDays, days)
	}
	return nil
}

var clusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "DnsName", JSONPath: "metadata.dnsName", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Instances", JSONPath: "properties.instances.count", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Cores", JSONPath: "properties.instances.cores"},
	{Name: "Ram", JSONPath: "properties.instances.ram"},
	{Name: "StorageSize", JSONPath: "properties.instances.storageSize"},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "BackupLocation", JSONPath: "properties.backup.location"},
	{Name: "RetentionDays", JSONPath: "properties.backup.retentionDays"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "LogsEnabled", JSONPath: "properties.logsEnabled"},
	{Name: "MetricsEnabled", JSONPath: "properties.metricsEnabled"},
	{Name: "DatacenterId", JSONPath: "properties.connection.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connection.lanId"},
	{Name: "Cidr", JSONPath: "properties.connection.primaryInstanceAddress"},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "MariaDB Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb-v2 cluster` allow you to manage the MariaDB Clusters under your account.",
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
	clusterCmd.AddCommand(ClusterGetCmd())

	return clusterCmd
}
