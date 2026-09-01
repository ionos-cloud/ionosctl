package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "snapshot",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List In-Memory DB Snapshots",
		LongDesc:  "List the automatic point-in-time Snapshots of your In-Memory DB Replica Sets. The Time column shows when each snapshot was dumped, and ReplicasetId / DatacenterId show which replica set and datacenter it belongs to. Queries every location by default; pin one with --location.",
		Example:   "ionosctl dbaas in-memory-db snapshot list",
		PreCmdRun: func(c *core.PreCommandConfig) error {

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				apiClient := inmemorydb.NewAPIClient(cfg)
				ls, _, err := apiClient.SnapshotApi.
					SnapshotsGet(context.Background()).Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
