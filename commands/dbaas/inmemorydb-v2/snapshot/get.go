package snapshot

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func SnapshotGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-inmemorydb-v2",
		Resource:   "snapshot",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an In-Memory DB Snapshot",
		Example:    "ionosctl dbaas in-memory-db-v2 snapshot get --snapshot-id <snapshot-id>",
		LongDesc:   "Use this command to retrieve details about an In-Memory DB Snapshot by using its ID.\n\nRequired values to run command:\n\n* Snapshot Id",
		PreCmdRun:  PreRunSnapshotId,
		CmdRun:     RunSnapshotGet,
		InitClient: true,
	})

	get.AddUUIDFlag(constants.FlagSnapshotId, constants.FlagIdShort, "", "The unique ID of the Snapshot", core.RequiredFlagOption(),
		core.WithCompletion(completer.SnapshotIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	return get
}

func PreRunSnapshotId(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagSnapshotId)
}

func RunSnapshotGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))
	c.Verbose("Getting Snapshot %v...", snapshotId)

	snapshot, _, err := client.Must().InMemoryDBClientV2.SnapshotsApi.SnapshotsFindById(
		context.Background(), snapshotId).Execute()
	if err != nil {
		return fmt.Errorf("could not get snapshot: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(snapshotCols, snapshot, cols))
}
