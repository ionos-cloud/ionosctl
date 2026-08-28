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

A replica set is the running database, either a single standalone instance or a leader-follower replication with one active and n-1 passive replicas; snapshots are read-only, point-in-time dumps taken automatically that can restore a replica set or seed a new one.

In-Memory DB is regional: every command targets a specific location. Set it with --location (or the IONOS_API_URL / config-file override); list commands query all locations by default.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(replicaset.Root())
	cmd.AddCommand(snapshot.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.InMemoryDB, "in-memory-db"}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations)
}
