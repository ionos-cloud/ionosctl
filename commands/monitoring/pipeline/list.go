package pipeline

import (
	"context"

	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func MonitoringListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List monitoring pipelines",
		LongDesc:  `List monitoring pipelines with their IDs, names, endpoints, and status. Because pipelines are regional, this command queries every supported location and merges the results by default; pass --location to restrict it to a single region. The output does not include ingest keys.`,
		Example: `# List pipelines across all regions
ionosctl monitoring pipeline list

# List pipelines in a single region, newest first
ionosctl monitoring pipeline list --location de/txl --order-by -createdDate`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				monitoringClient := monitoring.NewAPIClient(cfg)
				req := monitoringClient.PipelinesApi.PipelinesGet(context.Background())

				if fn := core.GetFlagName(c.NS, constants.FlagOrderBy); viper.IsSet(fn) {
					req = req.OrderBy(viper.GetString(fn))
				}

				ls, _, err := req.Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
