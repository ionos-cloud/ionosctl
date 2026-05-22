package central

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	cmd.AddCommand(CentralFindByIdCmd())
	cmd.AddCommand(CentralDisable())
	cmd.AddCommand(CentralEnable())

	return cmd
}

func enable(c *core.CommandConfig, enabled bool) error {

	logClient := logging.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
	input := logging.NewCentralLogging(enabled)

	a, _, err := logClient.CentralApi.CentralGet(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("failed getting the Id: %w", err)
	}

	r, _, err := logClient.CentralApi.CentralPut(context.Background(), a.Items[0].Id).
		CentralLoggingEnsure(logging.CentralLoggingEnsure{
			Properties: *input,
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed enabling the central: %w", err)
	}

	return c.Printer(allCols).Print(r)
}
