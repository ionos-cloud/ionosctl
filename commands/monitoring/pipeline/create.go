package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

func MonitoringPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a monitoring pipeline",
		LongDesc: `Create a new monitoring pipeline in the region given by --location. The only configurable property at creation is the name; the service provisions the rest and returns:
  * HttpEndpoint    - Prometheus remote-write target. Push metrics to <HttpEndpoint>/api/v1/push with the ingest key in the APIKEY header.
  * GrafanaEndpoint - the managed Grafana base URL for querying the ingested metrics.
  * Status          - the provisioning state; the pipeline is usable once it reports as available.

The ingest key is NOT part of this command's output. It is shown only once, immediately after creation, so retrieve or rotate it with 'monitoring key create --pipeline-id <id>' before configuring agents. An account may hold up to 10 pipelines by default (adjustable via Support).`,
		Example: `# Create a pipeline in Berlin
ionosctl monitoring pipeline create --location de/txl --name my-metrics

# Create a pipeline in Frankfurt and immediately generate its ingest key, capturing both
PIPE=$(ionosctl monitoring pipeline create --location de/fra --name prod-metrics --cols Id --no-headers)
ionosctl monitoring key create --location de/fra --pipeline-id "$PIPE"`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := monitoring.Pipeline{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			//later when put is fixed
			//z, _, err := client.Must().Monitoring.PipelinesApi.PipelinesPut(context.Background(), uuidgen.Must()).
			//	PipelineEnsure(monitoring.PipelineEnsure{
			//		Properties: input,
			//	}).Execute()
			z, _, err := client.Must().Monitoring.PipelinesApi.PipelinesPost(context.Background()).
				PipelineCreate(monitoring.PipelineCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating pipeline: %w", err)
			}

			return c.Printer(allCols).Print(z)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "A human-friendly label for the pipeline, shown in listings and the DCD. It does not affect the ingest or Grafana endpoints")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
