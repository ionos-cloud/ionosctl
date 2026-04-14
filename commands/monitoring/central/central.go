package central

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Enabled", JSONPath: "properties.enabled", Default: true},
	{Name: "GrafanaEndpoint", JSONPath: "metadata.grafanaEndpoint", Default: true},
	{Name: "Products", JSONPath: "metadata.products.*", Default: true},
}

func CentralCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "central",
			Aliases:          []string{"c"},
			Short:            "Central monitoring is a feature that allows you to use other products to send metrics to a central location.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(CentralFindByIdCmd())
	cmd.AddCommand(CentralDisable())
	cmd.AddCommand(CentralEnable())

	return cmd
}

func enable(c *core.CommandConfig, enabled bool) error {

	input := monitoring.CentralMonitoring{Enabled: enabled}

	r, _, err := client.Must().Monitoring.CentralApi.CentralPut(context.Background(), "").
		CentralMonitoringEnsure(monitoring.CentralMonitoringEnsure{
			Properties: input,
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed changing the enabled state: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(allCols, r, cols))
}
