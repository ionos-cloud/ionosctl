package logging_service

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/central"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "logging-service",
			Aliases: []string{"log-svc"},
			Short:   "Manage IONOS Logging Service (LaaS) pipelines",
			Long: `The Logging Service (LaaS) collects logs from your workloads and ships them to a managed Loki/Grafana backend so you can search and visualise them centrally.

The domain has three parts:
  - pipeline: the top-level resource. A pipeline exposes ingestion endpoints (a TCP address and an HTTP address) and a Grafana address for querying. Each pipeline is regional, so most commands require --location. Authentication to the ingestion endpoints uses a pipeline key (see 'pipeline key').
  - logs: the individual log streams inside a pipeline. Each log entry pairs a source (docker, systemd, kubernetes, generic) and a tag with a shipping protocol (http or tcp) and a destination (Loki, with a retention period). A pipeline must always contain at least one log.
  - central: central logging lets other IONOS products forward their logs into one shared Grafana view for the contract.

Typical flow: create a pipeline (with its first log) -> generate a key -> point your log shipper at the pipeline's TCP/HTTP address using that key -> query in Grafana.`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(pipeline.PipelineCmd())
	cmd.AddCommand(logs.LogsCmd())
	cmd.AddCommand(central.CentralCommand())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Logging}, constants.LoggingApiRegionalURL, constants.LoggingLocations)
}
