package inmemorydb_v2

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "in-memory-db-v2",
			Aliases:          []string{"inmemorydb-v2", "memdb-v2", "imdb-v2", "in-mem-db-v2", "inmemdb-v2"},
			Short:            "DBaaS In-Memory-DB V2 Operations",
			Long:             "The sub-commands of `ionosctl dbaas in-memory-db-v2` allow you to perform operations on In-Memory-DB V2 resources.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(cluster.ClusterCmd())
	cmd.AddCommand(snapshot.SnapshotCmd())
	cmd.AddCommand(version.VersionCmd())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.InMemoryDB, "in-memory-db"}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations)
}
