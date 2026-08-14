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
		Example:   "ionosctl dbaas inmemorydb snapshot list",
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
