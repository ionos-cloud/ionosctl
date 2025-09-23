package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

func MonitoringPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a pipeline's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl monitoring pipeline update --location de/txl --pipeline-id ID --name name",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagPipelineID); err != nil {
				return err
			}
			return nil
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

	cmd.AddStringFlag(constants.FlagPipelineID, "", "", constants.DescMonitoringPipeline, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.PipelineIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The new name of the Monitoring Pipeline", core.RequiredFlagOption())

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

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.MonitoringPipeline, rn,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
