package inmemorydb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/replicaset"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/snapshot"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "in-memory-db",
			Aliases:          []string{"inmemorydb", "memdb", "imdb", "in-mem-db", "inmemdb"},
			Short:            "DBaaS In-Memory-DB Operations",
			Long:             "The sub-commands of `ionosctl dbaas in-memory-db` allow you to perform operations on In-Memory-DB resources.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(replicaset.Root())
	cmd.AddCommand(snapshot.Root())

	return core.WithRegionalConfigOverride(cmd, []string{"in-memory-db"}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations)
}
