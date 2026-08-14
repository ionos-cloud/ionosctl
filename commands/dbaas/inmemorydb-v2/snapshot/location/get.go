package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func LocationGetCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-inmemorydb-v2",
		Resource:   "location",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an In-Memory DB Snapshot Location",
		LongDesc:   "Use this command to retrieve details about an In-Memory DB Snapshot Location by using its ID.\n\nRequired values to run command:\n\n* Snapshot Location Id",
		Example:    "ionosctl dbaas in-memory-db-v2 snapshot location get --snapshot-location-id <snapshot-location-id>",
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationGet,
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagSnapshotLocationId, constants.FlagIdShort, "", "The ID of the In-Memory DB Snapshot Location", core.RequiredFlagOption(),
		core.WithCompletion(completer.SnapshotLocationIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	return cmd
}

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagSnapshotLocationId)
}

func RunLocationGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	locationId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotLocationId))

	c.Verbose("Getting Snapshot Location...")

	location, _, err := client.Must().InMemoryDBClientV2.SnapshotLocationsApi.SnapshotlocationsFindById(context.Background(), locationId).Execute()
	if err != nil {
		return err
	}

	return c.Out(table.Sprint(locationCols, location, cols))
}
