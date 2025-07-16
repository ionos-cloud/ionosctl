package completer

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
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

	status := "disabled"
	if r.Properties.Enabled {
		status = "enabled"
	}

	_, err = fmt.Fprintf(c.Command.Command.OutOrStdout(), "The CentralMonitoring with id: %s is %s\n", r.Id, status)
	if err != nil {
		return fmt.Errorf("failed displaying the output: %w", err)
	}
	return nil
}
