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
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "restore",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create an In-Memory DB Restore",
		Example:   "ionosctl dbaas inmemorydb restore create " + core.FlagsUsage(constants.FlagReplicasetID, constants.FlagSnapshotId, constants.FlagName, constants.FlagDescription),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagReplicasetID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := inmemorydb.Restore{}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagReplicasetID)) {
				input.ReplicasetId = viper.GetString(core.GetFlagName(c.NS, constants.FlagReplicasetID))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
				input.DisplayName = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
				input.Description = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription)))
			}

			restore, _, err := client.Must().InMemoryDBClient.RestoreApi.SnapshotsRestoresPost(
				context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))).
				RestoreCreate(inmemorydb.RestoreCreate{Properties: input}).Execute()

			if err != nil {
				return fmt.Errorf("error creating restore: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasInMemoryDBSnapshot,
				restore, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagSnapshotId, "", "",
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
		core.RequiredFlagOption(),
	)

	cmd.AddStringFlag(constants.FlagReplicasetID, "", "",
		"The ID of the replica set the restore was applied on",
		core.WithCompletion(utils.ReplicasetIDs, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
		core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The human readable name of your snapshot")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "A description of the snapshot")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
