package key

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func KeyPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "key",
		Verb:      "create",
		Aliases:   []string{"post", "c", "generate"},
		ShortDesc: "Generate (rotate) a pipeline's ingest key",
		LongDesc: `Generate a fresh ingest key for a pipeline and print it. This is the only way to recover a key after creation, since it is shown just once.

Rotating is destructive to the old key: the previous key is invalidated immediately, so update every agent pushing to this pipeline with the new value. Agents send the key as the APIKEY header when writing to <HttpEndpoint>/api/v1/push (Prometheus remote-write).

The command prints only the raw key so it can be captured directly, e.g. piped into a secret store.`,
		Example: `# Rotate the key and print it
ionosctl monitoring key create --location de/txl --pipeline-id PIPELINE_ID

# Capture the new key into a variable
KEY=$(ionosctl monitoring key create --location de/txl --pipeline-id PIPELINE_ID)`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagPipelineID)
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

			c.Msg("%s", smth.Key)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagPipelineID, constants.FlagIdShort, "", "The ID of the pipeline whose ingest key should be rotated (from 'pipeline list')",
		core.WithCompletion(func() []string {
			return completer.PipelineIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
