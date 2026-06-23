package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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
	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		return deleteAll(c)
	}

	if err := c.RequireExplicitLocation(); err != nil {
		return err
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

	c.Msg("Successfully deleted logging pipeline")

	return nil
}

func pipelineSummary(p logging.PipelineRead) string {
	t := table.New(allCols)
	if extractErr := t.Extract(p); extractErr == nil {
		if rows := t.Rows(); len(rows) > 0 {
			return strings.TrimSpace(fmt.Sprintf("%v (%v)", rows[0]["Id"], rows[0]["Name"]))
		}
	}
	return p.Id
}

func deleteAll(c *core.CommandConfig) error {
	// Gather pipelines from every location (unless --location pins one), tagging each with its
	// location and location-scoped client, then hand the flat list to core.DeleteAll for a
	// consistent preview / per-item confirm-skip / summary flow.
	type located struct {
		pipeline logging.PipelineRead
		loc      string
		api      *logging.APIClient
	}
	var items []located
	if err := c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		lc := logging.NewAPIClient(cfg)
		pipelines, _, err := lc.PipelinesApi.PipelinesGet(context.Background()).Execute()
		if err != nil {
			return fmt.Errorf("failed listing pipelines in %s: %w", location, err)
		}
		for _, p := range pipelines.GetItems() {
			items = append(items, located{pipeline: p, loc: location, api: lc})
		}
		return nil
	}); err != nil {
		return err
	}

	return core.DeleteAll(c, core.DeleteAllOptions[located]{
		Resource: "Logging-Service Pipeline",
		List:     func() ([]located, error) { return items, nil },
		Summary: func(l located) string {
			return fmt.Sprintf("%s (location: %s)", pipelineSummary(l.pipeline), l.loc)
		},
		ID: func(l located) string { return l.pipeline.Id },
		Delete: func(l located) error {
			_, delErr := l.api.PipelinesApi.PipelinesDelete(context.Background(), l.pipeline.Id).Execute()
			return delErr
		},
	})
}
