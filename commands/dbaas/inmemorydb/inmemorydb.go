package inmemorydb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/replicaset"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "in-memory-db",
			Aliases: []string{"inmemorydb", "memdb", "imdb", "in-mem-db", "inmemdb"},
			Short:   "DBaaS In-Memory DB Operations",
			Long: `Manage IONOS CLOUD DBaaS In-Memory DB, a fully managed, Redis-compatible in-memory data store.

The domain has two resources:
  - replicaset: the running database. A replica set is either a single standalone instance (1 replica) or a leader-follower replication with one active and n-1 passive replicas. It carries the version, per-instance resources (cores/RAM), persistence mode, eviction policy, network connection (datacenter/LAN/CIDR), credentials, and a weekly maintenance window.
  - snapshot: read-only, point-in-time dumps of a replica set, taken automatically. A snapshot lives in the same datacenter as its replica set and is not available in other datacenters. Its child 'restore' resource rolls an existing replica set back to a snapshot's state; a brand-new replica set can also be created from a snapshot via 'replicaset create --snapshot-id'.

In-Memory DB is regional: every command targets a specific location. Set it with --location (or the IONOS_API_URL / config-file override); list commands query all locations by default.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(replicaset.Root())
	cmd.AddCommand(snapshot.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.InMemoryDB, "in-memory-db"}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations)
}
