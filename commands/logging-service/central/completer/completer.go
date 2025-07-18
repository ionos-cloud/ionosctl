package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
)

func CentralIDs() []string {
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
