package logs

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func LogsListCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "List the log streams of logging pipelines",
			LongDesc:  `List log streams. Pass --pipeline-id (with --location) to list the logs of one pipeline, or --all to list logs across pipelines. With --all and no --location, every location is queried and each log is annotated with its parent PipelineId and Location.`,
			Example: `ionosctl logging-service logs list --location de/txl --pipeline-id ID
ionosctl logging-service logs list --all`,
			PreCmdRun: preRunListCmd,
			CmdRun:    runListCmd,
		},
	)
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List logs from all pipelines. When --location is unset, logs from all locations are listed")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the pipeline whose log streams to list", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func preRunListCmd(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagLoggingPipelineId},
	); err != nil {
		return err
	}

	// --all fans out over every location; a specific --pipeline-id is
	// location-scoped and cannot be inferred, so require --location for it.
	if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		return c.RequireExplicitLocation()
	}
	return nil
}

func runListCmd(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		return listAll(c)
	}

	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pipelineId,
	).Execute()
	if err != nil {
		return err
	}

	return handleLogPrint(pipeline, c)
}

func listAll(c *core.CommandConfig) error {
	// When --location is unset, ListAllLocations queries every location
	// concurrently, listing all pipelines (and their logs) per location and
	// merging the results with a Location column.
	return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
		lc := logging.NewAPIClient(cfg)
		pipelines, _, err := lc.PipelinesApi.PipelinesGet(context.Background()).Execute()
		if err != nil {
			return nil, err
		}
		return flattenPipelineLogs(pipelines), nil
	})
}
