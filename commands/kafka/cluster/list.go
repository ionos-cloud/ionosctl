package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "List Kafka clusters",
			LongDesc:  `List Kafka clusters. Without --location every region is queried and results are merged with a Location column; pass --location to scope to one region. Narrow further with --name and --state.`,
			Example: `ionosctl kafka cl list
ionosctl kafka cl list --location de/txl --state AVAILABLE`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
					client := kafka.NewAPIClient(cfg)

					req := client.ClustersApi.ClustersGet(context.Background())
					if fn := core.GetFlagName(c.NS, constants.FlagFilterState); viper.IsSet(fn) {
						req = req.FilterState(viper.GetString(fn))
					}
					if fn := core.GetFlagName(c.NS, constants.FlagFilterName); viper.IsSet(fn) {
						req = req.FilterName(viper.GetString(fn))
					}

					ls, _, err := req.Execute()
					if err != nil {
						return nil, fmt.Errorf("failed listing kafka clusters: %w", err)
					}

					return ls, nil
				})
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagFilterName, "", "", "Return only clusters whose name contains this substring")
	cmd.AddSetFlag(
		constants.FlagFilterState, "", "", []string{"AVAILABLE", "BUSY", "DEPLOYING", "UPDATING", "FAILED_UPDATING", "FAILED", "DESTROYING"},
		"Return only clusters in this provisioning state (AVAILABLE = ready to use, BUSY/DEPLOYING/UPDATING = in progress)",
	)

	return cmd
}
