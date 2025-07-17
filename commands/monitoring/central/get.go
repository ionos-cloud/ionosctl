package central

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/central/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func CentralFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "central",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a CentralMonitoring",
		Example:   "ionosctl monitoring central get --location de/txl --central-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagCentralID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			centralId := viper.GetString(core.GetFlagName(c.NS, constants.FlagCentralID))

			r, _, err := client.Must().Monitoring.CentralApi.CentralFindById(context.Background(), centralId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the CentralMonitoring with ID: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.MonitoringCentral, r,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCentralID, "", "", fmt.Sprintf("%s is required ", constants.DescMonitoringCentral),
		core.WithCompletion(func() []string {
			return completer.CenralIDs()
		}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
