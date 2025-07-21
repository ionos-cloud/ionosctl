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
			r, _, err := client.Must().LoggingServiceClient.CentralApi.CentralLoggingGet(context.Background()).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the CentralLogging: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.LoggingCentral, r,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			_, err = fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			if err != nil {
				return err
			}

			return nil
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
