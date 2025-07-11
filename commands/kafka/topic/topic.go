package topic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "Name", "ReplicationFactor", "NumberOfPartitions", "RetentionTime", "SegmentByes", "State"}
	allCols     = []string{
		"Id", "Name", "ReplicationFactor", "NumberOfPartitions", "RetentionTime", "SegmentByes", "ClusterId", "State",
	}
)

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "topic",
			Short:            "The sub-commands of 'ionosctl kafka topic' allow you to manage kafka topics",
			Aliases:          []string{"cl"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(defaultCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(listCmd())

	return cmd
}
