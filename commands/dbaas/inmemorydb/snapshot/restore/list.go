package restore

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "restore",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List the Restores of a Snapshot",
		LongDesc:  "List the restore operations that have been performed from a given Snapshot (--snapshot-id). Each row shows the target replica set, the restore's state and its RestoreTime.",
		Example:   "ionosctl dbaas in-memory-db snapshot restore list --snapshot-id SNAPSHOT_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagSnapshotId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().InMemoryDBClient.RestoreApi.
				SnapshotsRestoresGet(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagSnapshotId, constants.FlagIdShort, "",
		"The ID of the Snapshot whose restore operations you want to list", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				// for each snapshot
				return utils.SnapshotProperty(func(snapshot inmemorydb.SnapshotRead) string {
					// return its ID
					return snapshot.Id + "\t" + snapshot.Metadata.SnapshotTime.Format("2006-01-02 15:04:05")
				})
			}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations,
		),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
