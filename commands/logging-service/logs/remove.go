package logs

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LogsRemoveCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "remove",
			ShortDesc: "Remove a log from a logging pipeline",
			Example:   `ionosctl logging-service logs remove --pipeline-id ID --log-tag TAG`,
			PreCmdRun: nil,
			CmdRun:    nil,
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
		constants.FlagLoggingPipelineLogTag, "", "", "The tag of the pipeline log that you want to remove",
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
