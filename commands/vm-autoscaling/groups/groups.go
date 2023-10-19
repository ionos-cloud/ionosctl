package groups

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/spf13/cobra"
)

func GroupsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "groups",
			Aliases:          []string{"g", "group"},
			Short:            "Autoscaling Groups Operations",
			Long:             "The sub-commands of `ionosctl autoscaling groups` allow you to manage the Autoscaling Groups under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(GroupCreateCmd())

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"GroupId":      "id",
		"Name":         "properties.name",
		"MinReplicas":  "properties.minReplicaCount",
		"MaxReplicas":  "properties.maxReplicaCount",
		"DatacenterId": "properties.datacenter.id",
		"State":        "metadata.state",
	}

	allCols     = []string{"GroupId", "Name", "MinReplicas", "MaxReplicas", "DatacenterId", "State"}
	defaultCols = allCols
)
