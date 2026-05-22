package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	monitoring "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MonitoringFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a pipeline",
		Example:   "ionosctl monitoring pipeline get --location de/txl --pipeline-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagPipelineID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))

			monClient := monitoring.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
			r, _, err := monClient.PipelinesApi.PipelinesFindById(context.Background(), pipelineId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the pipeline with ID '%s': %w", pipelineId, err)
			}

			return c.Printer(allCols).Print(r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagPipelineID, constants.FlagIdShort, "", constants.DescMonitoringPipeline)
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagPipelineID, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PipelineIDs(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
