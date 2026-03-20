package pipeline

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "GrafanaAddress", JSONPath: "properties.grafanaAddress", Default: true},
	{Name: "TCPAddress", JSONPath: "properties.tcpAddress"},
	{Name: "HTTPAddress", JSONPath: "properties.httpAddress"},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

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

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(PipelineListCmd())
	cmd.AddCommand(PipelineGetCmd())
	cmd.AddCommand(PipelineDeleteCmd())
	cmd.AddCommand(PipelineCreateCmd())
	cmd.AddCommand(PipelineUpdateCmd())
	cmd.AddCommand(PipelineKeyCmd())
	return cmd
}

func handlePipelinePrint(p logging.PipelineRead, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, p, cols))
}

func handleProvisioningPipelinePrint(p logging.PipelineRead, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, p, cols))
}
