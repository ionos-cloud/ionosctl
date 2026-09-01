package replicaset

import (
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
			Use:     "replicaset",
			Aliases: []string{"rs", "replica-set", "replicasets", "cluster"},
			Short:   "Manage In-Memory DB Replica Sets (the running databases)",
			Long: `The sub-commands of ` + "`ionosctl dbaas in-memory-db replicaset`" + ` manage In-Memory DB Replica Sets, the Redis-compatible databases themselves.

A replica set runs in one of two modes, determined by the --replicas count:
  - Standalone (replicas = 1): a single instance, no redundancy.
  - Replication / leader-follower (replicas > 1): one active instance plus n-1 passive replicas. The passive replicas are hot standbys that take over if the active fails; they are NOT read replicas and cannot serve reads.

Every replica set has: an engine version, per-instance resources (cores + RAM; storage is derived automatically), a persistence mode and eviction policy, exactly one network connection (datacenter + LAN + CIDR), an initial user, and a weekly 4-hour maintenance window. A replica set can be created empty or restored from an existing snapshot (see 'replicaset create --snapshot-id').`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Create())
	cmd.AddCommand(Get())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	// cmd.AddCommand(Update()) // Update is disabled until an API fix is rolled out

	return cmd
}
