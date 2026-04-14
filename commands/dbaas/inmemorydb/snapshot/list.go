package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			ls, _, err := client.Must().InMemoryDBClient.SnapshotApi.
				SnapshotsGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
