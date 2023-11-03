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
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
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
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServicePipelineIds(), cobra.ShellCompDirectiveNoFileComp
		},
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
		return deleteAll(c)
	}

	pipelineId, err := c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineId)
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete %s", pipelineId)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, _, err = client.Must().LoggingServiceClient.PipelinesApi.PipelinesDelete(
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

	pipelinesConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.LoggingServicePipeline, pipelines)
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(
		*items, func(p ionoscloud.Pipeline) error {
			pInfo := completions.NewCompleter(pipelinesConverted, "Id").AddInfo("Name", "(%v)").ToString()[0]
			pInfo = strings.Replace(pInfo, "\t", "", 1)

			yes := confirm.FAsk(
				c.Command.Command.InOrStdin(), fmt.Sprintf(
					"delete %s", pInfo,
				),
				viper.GetBool(constants.ArgForce),
			)
			if yes {
				_, delErr := client.Must().DnsClient.ZonesApi.ZonesDelete(c.Context, *p.Id).Execute()
				if delErr != nil {
					return fmt.Errorf("failed deleting %s: %w", pInfo, delErr)
				}
			}
			return nil
		},
	)

	return err
}
