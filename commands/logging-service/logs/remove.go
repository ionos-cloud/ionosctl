package logs

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/viper"
)

func LogsRemoveCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "remove",
			ShortDesc: "Remove a log from a logging pipeline. NOTE:" +
				"There needs to be at least one log in a pipeline at all times.",
			Example:   `ionosctl logging-service logs remove --pipeline-id ID --log-tag TAG`,
			PreCmdRun: preRunRemoveCmd,
			CmdRun:    runRemoveCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL),
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "The tag of the pipeline log that you want to delete",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func runRemoveCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete %s", tag), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	var newLogs []ionoscloud.PipelineCreatePropertiesLogs
	for _, log := range *patchPipeline.Properties.Logs {
		if *log.Tag == tag {
			continue
		}

		newLogs = append(newLogs, log)
	}
	patchPipeline.Properties.Logs = &newLogs

	_, _, err = client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).Pipeline(
		*patchPipeline,
	).Execute()
	if err != nil {
		return err
	}

	fmt.Fprintf(
		c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(
			"Log successfully removed from pipeline",
		),
	)

	return nil
}

func preRunRemoveCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}
