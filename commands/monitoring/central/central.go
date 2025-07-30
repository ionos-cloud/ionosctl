package central

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Enabled", "GrafanaEndpoint", "Products"}
)

func CentralCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "central",
			Aliases:          []string{"c"},
			Short:            "Central monitoring is a feature that allows you to use other products to send metrics to a central location.",
			TraverseChildren: true,
		},
	}

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

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.MonitoringCentral, r,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return fmt.Errorf("failed generating the output: %w", err)
	}

	_, err = fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
