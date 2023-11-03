package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	logsvc "github.com/ionos-cloud/sdk-go-logging"
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
			CmdRun:    run,
		},
	)

	return cmd
}

func run(cmd *core.CommandConfig) error {
	cfg := logsvc.NewConfigurationFromEnv()
	cl := logsvc.NewAPIClient(cfg)
	_, _, err := cl.PipelinesApi.PipelinesGet(context.Background()).Execute()
	//_, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	// NOTE: why won't it work with Must() ????
	if err != nil {
		return err
	}
	//client.Must()

	fmt.Println("yey")
	return nil
}
