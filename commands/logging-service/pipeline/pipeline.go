package pipeline

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
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

func handlePipelinePrint(p logging.Pipeline, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.LoggingServicePipeline, p, tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func handleProvisioningPipelinePrint(p logging.PipelineRead, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.LoggingServicePipeline, p, tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
