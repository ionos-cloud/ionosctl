package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

func MonitoringPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "monitoring",
		Resource:  "pipeline",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create an pipeline",
		Example:   "ionosctl monitoring pipeline create --location de/txl --name name",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := monitoring.Pipeline{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			//later when put is fixed
			//z, _, err := client.Must().Monitoring.PipelinesApi.PipelinesPut(context.Background(), uuidgen.Must()).
			//	PipelineEnsure(monitoring.PipelineEnsure{
			//		Properties: input,
			//	}).Execute()
			z, _, err := client.Must().Monitoring.PipelinesApi.PipelinesPost(context.Background()).
				PipelineCreate(monitoring.PipelineCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating pipeline: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.MonitoringPipeline, z,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the JSON output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the Monitoring pipeline")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
