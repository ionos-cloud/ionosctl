package topic

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "ReplicationFactor", JSONPath: "properties.replicationFactor", Default: true},
	{Name: "NumberOfPartitions", JSONPath: "properties.numberOfPartitions", Default: true},
	{Name: "RetentionTime", JSONPath: "properties.logRetention.retentionTime", Default: true},
	{Name: "SegmentByes", JSONPath: "properties.logRetention.segmentBytes", Default: true},
	{Name: "ClusterId", JSONPath: "href"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func Command() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "topic",
			Short:            "The sub-commands of 'ionosctl kafka topic' allow you to manage kafka topics",
			Aliases:          []string{"t"},
			TraverseChildren: true,
		},
	}
	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(listCmd())

	return cmd
}
