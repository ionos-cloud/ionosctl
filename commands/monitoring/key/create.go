package key

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/viper"
)

func KeyPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "key",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a new key for a pipeline",
		Example:   "ionosctl monitoring key create --location de/txl --pipeline-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagPipelineID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))

			_, _, err := client.Must().Monitoring.PipelinesApi.PipelinesFindById(context.Background(), pipelineId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the pipeline with ID '%s': %w", pipelineId, err)
			}

			smth, _, err := client.Must().Monitoring.KeyApi.PipelinesKeyPost(context.Background(), pipelineId).
				Body(map[string]interface{}{}).Execute()
			if err != nil {
				return fmt.Errorf("failed updating the key %s: %w", pipelineId, err)
			}

			_, err = fmt.Fprintf(c.Command.Command.OutOrStdout(), "The new key is: %s", jsontabwriter.GenerateRawOutput(smth.Key))
			if err != nil {
				return fmt.Errorf("failed writing the key to output %s: %w", pipelineId, err)
			}

			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagPipelineID, constants.FlagIdShort, "", fmt.Sprintf("%s. Required or -%s", constants.DescMonitoringPipeline, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			return completer.PipelineIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
