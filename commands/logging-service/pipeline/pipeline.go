package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "Name", "GrafanaAddress", "CreatedDate", "State"}
	allCols     = []string{"Id", "Name", "GrafanaAddress", "TCPAddress", "HTTPAddress", "CreatedDate", "State"}
)

func PipelineCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "pipeline",
			Aliases: []string{"p", "pipelines"},
			Short: "The subcommands of `ionosctl logging-service pipeline` allow you to manage logging pipelines. " +
				"They are the backbone of a centralized logging system, " +
				"referring to an instance or configuration of the logging service you can create",
		},
	}

	cmd.AddCommand(PipelineListCmd())
	cmd.AddCommand(PipelineGetCmd())
	cmd.AddCommand(PipelineDeleteCmd())
	cmd.AddCommand(PipelineCreateCmd())
	cmd.AddCommand(PipelineUpdateCmd())
	cmd.AddCommand(PipelineKeyCmd())
	return cmd
}

func handlePipelinePrint(p ionoscloud.Pipeline, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.LoggingServicePipeline, p, tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func PipelineStateInterrogator(_ *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(context.Background(), objId).Execute()
	if err != nil {
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}
