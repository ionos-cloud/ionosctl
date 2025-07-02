package pipeline

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func MonitoringListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "Retrieve pipelines",
		Example:   "ionosctl monitoring pipeline list --location de-fra",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagMonitoringPipelineLocation); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			req := client.Must().Monitoring.PipelinesApi.PipelinesGet(context.Background())

			req.OrderBy("pipeline_id")

			if fn := core.GetFlagName(c.NS, constants.FlagOrderBy); viper.IsSet(fn) {
				req = req.OrderBy(viper.GetString(fn))
			}

			ls, _, err := req.Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.MonitoringPipeline, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagOrderBy, "", "", "The field to order the results by. If not provided, the results will be ordered by the default field.")

	cmd.AddStringFlag(constants.FlagMonitoringPipelineLocation, constants.FlagMonitoringLocationShort, "", constants.DescMonitoringPipelineLocation)
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMonitoringPipelineLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GatewaysIDs(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
