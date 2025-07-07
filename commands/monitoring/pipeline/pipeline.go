package pipeline

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "GrafanaEndpoint", "HttpEndpoint", "Status"}
)

func PipelineCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "pipeline",
			Aliases:          []string{"p", "pipe"},
			Short:            "A metric pipeline refers to an instance or configuration of the Monitoring Service ",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddCommand(MonitoringListCmd())
	cmd.AddCommand(MonitoringDeleteCmd())

	return cmd
}
