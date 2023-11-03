package logs

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LogsListCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "Retrieve logging pipeline logs",
			Example:   "ionosctl logging-service logs list --pipeline-id ID",
			CmdRun:    runListCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Use this flag to list all logging pipeline logs")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline you want to list logs for", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServicePipelineIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func runListCmd(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		return listAll(c)
	}

	return nil
}

func listAll(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	pipelines, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	var logs [][]ionoscloud.ResponsePipeline
	for _, p := range *pipelines.Items {
		logs = append(logs, *p.Properties.Logs)
	}

	logsConverted, err := resource2table.ConvertLoggingServicePipelineLogsToTable(pipelines)
	if err != nil {
		return fmt.Errorf("could not convert Logging Service Pipeline Logs to table format: %w", err)
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		logs, logsConverted, tabheaders.GetHeaders(
			allCols,
			defaultCols, cols,
		),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
