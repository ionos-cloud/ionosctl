package replicaset

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.displayName", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "DNSName", JSONPath: "metadata.dnsName", Default: true},
	{Name: "Replicas", JSONPath: "properties.replicas", Default: true},
	{Name: "Cores", JSONPath: "properties.resources.cores", Default: true},
	{Name: "RAM", JSONPath: "properties.resources.ram", Default: true},
	{Name: "StorageSize", JSONPath: "properties.resources.storage", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "BackupLocation", JSONPath: "properties.backupLocation"},
	{Name: "PersistenceMode", JSONPath: "properties.persistenceMode"},
	{Name: "EvictionPolicy", JSONPath: "properties.evictionPolicy"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId"},
	{Name: "Username", JSONPath: "properties.credentials.username"},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "replicaset",
			Aliases:          []string{"rs", "replica-set", "replicasets", "cluster"},
			Short:            "The sub-commands of 'ionosctl dbaas inmemorydb replicaset' allow you to manage In-Memory DB Replica Sets.",
			Long:             "In-Memory DB replica set with support for a single instance or a In-Memory DB replication in leader follower mode. The mode is determined by the number of replicas. One replica is standalone, everything else an In-Memory DB replication as leader follower mode with one active and n-1 passive replicas.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(Create())
	cmd.AddCommand(Get())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	// cmd.AddCommand(Update()) // Update is disabled until an API fix is rolled out

	return cmd
}
