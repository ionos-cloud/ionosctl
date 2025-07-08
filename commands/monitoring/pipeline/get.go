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
		Example:   "ionosctl monitoring pipeline get --pipeline-id GATEWAYID --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagPipelineID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineID))

			r, _, err := client.Must().Monitoring.PipelinesApi.PipelinesFindById(context.Background(), pipelineId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the pipeline with ID '%s': %w", pipelineId, err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.MonitoringPipeline, r, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
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
