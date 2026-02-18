package flowlog

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultFlowLogCols = []string{"FlowLogId", "Name", "Action", "Direction", "Bucket", "State"}
)

func FlowlogCmd() *core.Command {
	flowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"fl"},
			Short:            "FlowLog Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl flowlog` + "`" + ` allow you to create, list, get, delete FlowLogs on specific NICs.`,
			TraverseChildren: true,
		},
	}
	globalFlags := flowLogCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultFlowLogCols, tabheaders.ColsMessage(defaultFlowLogCols))
	_ = viper.BindPFlag(core.GetFlagName(flowLogCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	flowLogCmd.AddCommand(FlowLogListCmd())
	flowLogCmd.AddCommand(FlowLogGetCmd())
	flowLogCmd.AddCommand(FlowLogCreateCmd())
	flowLogCmd.AddCommand(FlowLogDeleteCmd())

	return core.WithConfigOverride(flowLogCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
