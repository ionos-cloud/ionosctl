package central

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func CentralFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a CentralMonitoring",
		Example:   "ionosctl monitoring central get --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			centralId := viper.GetString(core.GetFlagName(c.NS, constants.FlagCentralID))

			r, _, err := client.Must().Monitoring.CentralApi.CentralFindById(context.Background(), centralId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the CentralMonitoring with ID: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "The CentralMonitoring with id: %s is %v\n", r.Id, r.Properties.Enabled)
			return nil
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
