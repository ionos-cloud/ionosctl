package pipeline

import (
	"context"

	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func PipelineListCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "Retrieve logging pipelines",
			Example:   "ionosctl logging-service pipeline list",
			CmdRun: func(c *core.CommandConfig) error {
				return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
					loggingClient := logging.NewAPIClient(cfg)
					ls, _, err := loggingClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
					return ls, err
				})
			},
			InitClient: true,
		},
	)

	return cmd
}
