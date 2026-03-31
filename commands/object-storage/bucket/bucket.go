package bucket

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "CreationDate", JSONPath: "CreationDate", Default: true},
	{Name: "Region", JSONPath: "Region", Default: true},
}

func BucketCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "bucket",
			Aliases:          []string{"b"},
			Short:            "Bucket operations for contract-owned object storage",
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

	cmd.AddCommand(ListBucketsCmd())
	cmd.AddCommand(CreateBucketCmd())
	cmd.AddCommand(GetBucketCmd())
	cmd.AddCommand(HeadBucketCmd())
	cmd.AddCommand(GetBucketVersioningCmd())
	cmd.AddCommand(DeleteBucketCmd())
	return cmd
}
