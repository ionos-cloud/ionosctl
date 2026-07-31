package replicaset

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List In-Memory DB Replica Sets",
		LongDesc:  "List In-Memory DB Replica Sets. By default this queries every location and merges the results; pin a single region with --location. Filter by name with --name.",
		Example:   "ionosctl dbaas in-memory-db replicaset list",
		PreCmdRun: func(c *core.PreCommandConfig) error {

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				apiClient := inmemorydb.NewAPIClient(cfg)
				ls, _, err := apiClient.ReplicaSetApi.
					ReplicasetsGet(context.Background()).Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Return only Replica Sets whose display name matches this value",
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
