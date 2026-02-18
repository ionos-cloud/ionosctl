package request

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultRequestCols = []string{"RequestId", "CreatedDate", "Method", "Status", "Message", "Targets"}
	allRequestCols     = []string{"RequestId", "CreatedDate", "CreatedBy", "Method", "Status", "Message", "Url", "Body", "Targets"}
)

func RequestCmd() *core.Command {
	reqCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "request",
			Aliases:          []string{"req"},
			Short:            "Request Operations",
			Long:             "The sub-commands of `ionosctl request` allow you to see information about requests on your account. With the `ionosctl request` command, you can list, get or wait for a Request.",
			TraverseChildren: true,
		},
	}
	globalFlags := reqCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultRequestCols, tabheaders.ColsMessage(allRequestCols))
	_ = viper.BindPFlag(core.GetFlagName(reqCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = reqCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allRequestCols, cobra.ShellCompDirectiveNoFileComp
	})

	reqCmd.AddCommand(RequestListCmd())
	reqCmd.AddCommand(RequestGetCmd())
	reqCmd.AddCommand(RequestWaitCmd())

	return core.WithConfigOverride(reqCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
