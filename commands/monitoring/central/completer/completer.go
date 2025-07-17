package completer

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

var (
	allCols = []string{"Id", "Enabled", "GrafanaEndpoint", "Products"}
)

func CentralEnable(c *core.CommandConfig, enabled bool) error {

	centralId := viper.GetString(core.GetFlagName(c.NS, constants.FlagCentralID))
	input := monitoring.CentralMonitoring{Enabled: enabled}

	r, _, err := client.Must().Monitoring.CentralApi.CentralPut(context.Background(), centralId).
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func CenralIDs() []string {
	central, _, err := client.Must().Monitoring.CentralApi.CentralGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	centralConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.MonitoringCentral, central)
	if err != nil {
		return nil
	}
	return completions.NewCompleter(centralConverted, "Id").AddInfo("Enabled").AddInfo("GrafanaEndpoint").AddInfo("Products").ToString()
}
