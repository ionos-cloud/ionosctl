package publicaccessblock

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "BlockPublicAcls", JSONPath: "BlockPublicAcls", Default: true},
	{Name: "IgnorePublicAcls", JSONPath: "IgnorePublicAcls", Default: true},
	{Name: "BlockPublicPolicy", JSONPath: "BlockPublicPolicy", Default: true},
	{Name: "RestrictPublicBuckets", JSONPath: "RestrictPublicBuckets", Default: true},
}

type publicAccessBlockInfo struct {
	BlockPublicAcls       bool `json:"BlockPublicAcls"`
	IgnorePublicAcls      bool `json:"IgnorePublicAcls"`
	BlockPublicPolicy     bool `json:"BlockPublicPolicy"`
	RestrictPublicBuckets bool `json:"RestrictPublicBuckets"`
}

func PublicAccessBlockCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "public-access-block",
			Aliases:          []string{"pab"},
			Short:            "Public access block operations for contract-owned object storage buckets",
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

	cmd.AddCommand(GetCmd())
	cmd.AddCommand(PutCmd())
	cmd.AddCommand(DeleteCmd())

	return cmd
}
