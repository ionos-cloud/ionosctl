package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func PipelineDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "delete",
			ShortDesc: "Delete a logging pipeline using its ID",
			Example:   "ionosctl logging-service pipeline delete --pipeline-id ID",
			PreCmdRun: preRunDeleteCmd,
			CmdRun:    runDeleteCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Use this flag to delete all logging pipelines")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline you want to delete", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func preRunDeleteCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.FlagLoggingPipelineId}, []string{constants.ArgAll},
	)
}

func runDeleteCmd(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		if err := deleteAll(c); err != nil {
			return err
		}

		fmt.Fprintf(
			c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(
				"Successfully deleted logging pipeline",
			),
		)

		return nil
	}

	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))

	if !confirm.FAsk(
		c.Command.Command.InOrStdin(), fmt.Sprintf("delete %s", pipelineId), viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesDelete(
		context.Background(), pipelineId,
	).Execute()
	if err != nil {
		return err
	}

	fmt.Fprintf(
		c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(
			"Successfully deleted logging pipeline",
		),
	)

	return nil
}

func deleteAll(c *core.CommandConfig) error {
	pipelines, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	items, ok := pipelines.GetItemsOk()
	if !ok || items == nil {
		return fmt.Errorf("could not retrieve Logging-Service Pipelines items")
	}

	if len(items) <= 0 {
		return fmt.Errorf("no Logging-Service Pipelines to delete")
	}

	err = functional.ApplyAndAggregateErrors(
		items, func(p logging.Pipeline) error {
			pipelineConverted, err := json2table.ConvertJSONToTable("", jsonpaths.LoggingServicePipeline, p)
			if err != nil {
				return err
			}

			pInfo := completions.NewCompleter(pipelineConverted, "Id").AddInfo("Name", "(%v)").ToString()[0]
			pInfo = strings.Replace(pInfo, "\t", "", 1)

			yes := confirm.FAsk(
				c.Command.Command.InOrStdin(), fmt.Sprintf(
					"delete %s", pInfo,
				),
				viper.GetBool(constants.ArgForce),
			)
			if yes {
				_, delErr := client.Must().LoggingServiceClient.PipelinesApi.PipelinesDelete(
					c.Context,
					*p.Id,
				).Execute()
				if delErr != nil {
					return fmt.Errorf("failed deleting %s: %w", pInfo, delErr)
				}
			}
			return nil
		},
	)

	return err
}
