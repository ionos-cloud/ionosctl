package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
			CmdRun:    runListCmd,
		},
	)

	return cmd
}

func runListCmd(c *core.CommandConfig) error {
	pipelines, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(pipelines)
}
