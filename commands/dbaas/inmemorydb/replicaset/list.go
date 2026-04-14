package replicaset

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List In-Memory DB Replica Sets",
		Example:   "ionosctl dbaas inmemorydb replicaset list",
		PreCmdRun: func(c *core.PreCommandConfig) error {

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().InMemoryDBClient.ReplicaSetApi.
				ReplicasetsGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "You can filter the Replica Sets by name",
		core.WithCompletion(
			func() []string {
				// for each replica set
				return utils.ReplicasetProperty(func(replica inmemorydb.ReplicaSetRead) string {
					// return its name
					return replica.Properties.DisplayName
				})
			},
			constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations,
		),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
