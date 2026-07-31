package tagging

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "Key", JSONPath: "Key", Default: true},
	{Name: "Value", JSONPath: "Value", Default: true},
}

type tagInfo struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func TaggingCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "tagging",
			Aliases: []string{"tag"},
			Short:   "Manage key/value tags on a bucket",
			Long: `Manage a bucket's tag set: key/value pairs attached to the bucket for organization, cost allocation and automation (e.g. Environment=production, Team=platform). Tags describe the bucket itself and are independent of object-level tags.

The bucket has one tag set; 'put' replaces it entirely (it is not an additive merge - include every tag you want to keep), 'get' lists the current tags, and 'delete' removes all tags from the bucket.`,
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
