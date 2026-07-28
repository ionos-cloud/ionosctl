package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func LocationListCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-inmemorydb-v2",
		Resource:   "location",
		Verb:       "list",
		Aliases:    []string{"ls"},
		ShortDesc:  "List In-Memory DB Snapshot Locations",
		LongDesc:   "Use this command to retrieve a list of Object Storage locations where snapshots can be stored.",
		Example:    "ionosctl dbaas in-memory-db-v2 snapshot location list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunLocationList,
		InitClient: true,
	})
	cmd.AddInt32Flag(constants.FlagLimit, "", 100, "The maximum number of elements to return")
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "The first element to return")
	return cmd
}

func RunLocationList(c *core.CommandConfig) error {
	req := client.Must().InMemoryDBClientV2.SnapshotLocationsApi.SnapshotlocationsGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagLimit); viper.IsSet(fn) {
		req = req.Limit(viper.GetInt32(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
		req = req.Offset(viper.GetInt32(fn))
	}
	locations, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(locationCols, locations, cols, table.WithPrefix("items")))
}
