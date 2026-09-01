package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

func MonitoringPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Rename a monitoring pipeline",
		LongDesc: `Update a pipeline's mutable properties. Today the only editable property is the name; the HTTP endpoint, Grafana endpoint, and ingest key are fixed for the life of the pipeline and are not affected here (rotate the key with 'monitoring key create').

Under the hood this reads the current pipeline and writes it back with the new values (GET + PUT), emulating a partial update, so unspecified properties are preserved.`,
		Example: `# Rename a pipeline
ionosctl monitoring pipeline update --location de/txl --pipeline-id PIPELINE_ID --name new-name`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagPipelineID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))
			g, _, err := client.Must().Monitoring.PipelinesApi.PipelinesFindById(context.Background(), pipelineId).Execute()
			if err != nil {
				return fmt.Errorf("failed retrieving pipeline with ID '%s': %w", pipelineId, err)
			}
			return partiallyUpdatePipelinePrint(c, g)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagPipelineID, "", "", "The ID of the monitoring pipeline to update (from 'pipeline list')", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.PipelineIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new human-friendly label for the pipeline. It does not affect the ingest or Grafana endpoints", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func partiallyUpdatePipelinePrint(c *core.CommandConfig, r monitoring.PipelineRead) error {
	input := r.Properties
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}

	pipelineid := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))
	rn, _, err := client.Must().Monitoring.PipelinesApi.PipelinesPut(context.Background(), pipelineid).
		PipelineEnsure(monitoring.PipelineEnsure{
			Properties: input,
		}).Execute()

	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(rn)
}
