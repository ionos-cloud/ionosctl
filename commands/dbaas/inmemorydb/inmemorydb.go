package inmemorydb

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "in-memory-db",
			Aliases:          []string{"inmemorydb", "memdb", "imdb", "inmemdb"},
			Short:            "DBaaS In-Memory-DB Operations",
			Long:             ``, // TODO
			TraverseChildren: true,
		},
	}

	_ = inmemorydb.Version

	return cmd
}
