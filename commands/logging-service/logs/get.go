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

func LogsGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "get",
			ShortDesc: "Retrieve a log from a logging pipeline",
			Example:   `ionosctl logging-service logs get --pipeline-id ID --log-tag TAG`,
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineId,
		func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServicePipelineIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "The tag of the pipeline log that you want to retrieve",
		core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineLogTag,
		func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServiceLogTags(
				viper.GetString(
					core.GetFlagName(
						cmd.NS, constants.FlagLoggingPipelineId,
					),
				),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func runGetCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	var log ionoscloud.PipelineResponse

	for _, l := range *pipeline.Properties.Logs {
		if *l.Tag == tag {
			log = l

			break
		}
	}

	logConverted, err := resource2table.ConvertLoggingServicePipelineLogToTable(log)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(
		log, logConverted, tabheaders.GetHeaders(
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

func preRunGetCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}
