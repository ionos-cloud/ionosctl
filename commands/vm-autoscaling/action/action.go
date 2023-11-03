package action

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "action",
			Aliases:          []string{"act"},
			Short:            "Autoscaling Servers Operations",
			Long:             "The sub-commands of `ionosctl autoscaling server` allow you to manage the Autoscaling Servers under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"ActionId": "id",
	}

	allCols = []string{
		"ActionId",
	}

	defaultCols = allCols
)
