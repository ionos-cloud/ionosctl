package cors

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "AllowedOrigins", JSONPath: "AllowedOrigins", Default: true},
	{Name: "AllowedMethods", JSONPath: "AllowedMethods", Default: true},
	{Name: "AllowedHeaders", JSONPath: "AllowedHeaders", Default: true},
	{Name: "ExposeHeaders", JSONPath: "ExposeHeaders"},
	{Name: "MaxAgeSeconds", JSONPath: "MaxAgeSeconds"},
	{Name: "ID", JSONPath: "ID"},
}

type corsRuleInfo struct {
	AllowedOrigins string `json:"AllowedOrigins"`
	AllowedMethods string `json:"AllowedMethods"`
	AllowedHeaders string `json:"AllowedHeaders"`
	ExposeHeaders  string `json:"ExposeHeaders"`
	MaxAgeSeconds  string `json:"MaxAgeSeconds"`
	ID             string `json:"ID"`
}

func CorsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cors",
			Short:            "Bucket CORS operations for contract-owned object storage",
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
