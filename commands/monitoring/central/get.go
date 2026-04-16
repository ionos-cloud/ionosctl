package central

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve CentralMonitoring",
		Example:   "ionosctl monitoring central get --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			r, _, err := client.Must().Monitoring.CentralApi.CentralGet(context.Background()).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the CentralMonitoring: %w", err)
			}

			return c.Printer(allCols).Prefix("items").Print(r)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
