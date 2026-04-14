package replicaset

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get an In-Memory DB Replica Set",
		Example:   fmt.Sprintf("ionosctl dbaas inmemorydb replicaset get %s", core.FlagsUsage(constants.FlagReplicasetID)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(
				c.Command, c.NS,
				[]string{constants.FlagReplicasetID},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagReplicasetID))

			rs, _, err := client.Must().InMemoryDBClient.ReplicaSetApi.ReplicasetsFindById(context.Background(), id).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(rs)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagReplicasetID, constants.FlagIdShort, "",
		"The ID of the Replica Set you want to delete",
		core.WithCompletion(utils.ReplicasetIDs, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
