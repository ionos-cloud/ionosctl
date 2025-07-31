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
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
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
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "The tag of the pipeline log that you want to retrieve",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
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

	var log logging.PipelineNoAddrLogs

	for _, l := range pipeline.Properties.Logs {
		if l.Tag == tag {
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func preRunGetCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}
