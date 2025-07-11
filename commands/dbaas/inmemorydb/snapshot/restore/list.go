package restore

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "restore",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List In-Memory DB Restores",
		Example:   "ionosctl dbaas inmemorydb restore list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagSnapshotId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().InMemoryDBClient.RestoreApi.
				SnapshotsRestoresGet(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DbaasInMemoryDBSnapshot, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagSnapshotId, constants.FlagIdShort, "",
		"The ID of the In-Memory DB Snapshot to list restore points from", core.RequiredFlagOption(),
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
