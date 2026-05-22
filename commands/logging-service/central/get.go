package central

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func CentralFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "logging-service",
		Resource:  "central",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve CentralLogging",
		Example:   "ionosctl logging-service central get --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			logClient := logging.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
			r, _, err := logClient.CentralApi.CentralGet(context.Background()).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the CentralLogging: %w", err)
			}

			return c.Printer(allCols).Prefix("items").Print(r)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
