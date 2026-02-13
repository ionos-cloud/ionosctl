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

			if c.Command.Command.Flags().Changed(constants.FlagReplicasetID) {
				replicasetId, err := c.Command.Command.Flags().GetString(constants.FlagReplicasetID)
				if err != nil {
					return err
				}
				input.ReplicasetId = replicasetId
			}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				input.DisplayName = pointer.From(name)
			}

			if c.Command.Command.Flags().Changed(constants.FlagDescription) {
				description, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
				if err != nil {
					return err
				}
				input.Description = pointer.From(description)
			}

			snapshotId, err := c.Command.Command.Flags().GetString(constants.FlagSnapshotId)
			if err != nil {
				return err
			}

			restore, _, err := client.Must().InMemoryDBClient.RestoreApi.SnapshotsRestoresPost(
				context.Background(), snapshotId).
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
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
