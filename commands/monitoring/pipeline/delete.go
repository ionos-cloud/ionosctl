package pipeline

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

func MonitoringDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a pipeline",
		Example:   "ionosctl monitoring pipeline delete --location de/txl --pipeline-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagPipelineID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			monitoringId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))
			z, _, err := client.Must().Monitoring.PipelinesApi.PipelinesFindById(context.Background(), monitoringId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting pipeline by id %s: %w", monitoringId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(),
				fmt.Sprintf("Are you sure you want to delete pipeline with name: %s, id: %s ", z.Properties.Name, z.Id),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, newerr := client.Must().Monitoring.PipelinesApi.PipelinesDelete(context.Background(), monitoringId).Execute()
			if newerr != nil {
				return fmt.Errorf("failed deleting the pipeline: %w", newerr)
			}

			return newerr
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagPipelineID, constants.FlagIdShort, "", fmt.Sprintf("%s. Required or -%s", constants.DescMonitoringPipeline, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			return completer.PipelineIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all pipelines."))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all pipelines!"))
	xs, _, err := client.Must().Monitoring.PipelinesApi.PipelinesGet(context.Background()).Execute()

	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(z monitoring.PipelineRead) error {
		yes :=
			confirm.FAsk(c.Command.Command.InOrStdin(),
				fmt.Sprintf("Are you sure you want to delete pipeline with name: %s, id: %s ", z.Properties.Name, z.Id),
				viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().Monitoring.PipelinesApi.PipelinesDelete(context.Background(), z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", z.Id, z.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
