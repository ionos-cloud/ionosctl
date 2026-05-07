package lifecycle

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var allCols = []table.Column{
	{Name: "ID", JSONPath: "ID", Default: true},
	{Name: "Prefix", JSONPath: "Prefix", Default: true},
	{Name: "Status", JSONPath: "Status", Default: true},
	{Name: "ExpirationDays", JSONPath: "ExpirationDays", Default: true},
	{Name: "ExpirationDate", JSONPath: "ExpirationDate"},
	{Name: "ExpiredObjectDeleteMarker", JSONPath: "ExpiredObjectDeleteMarker"},
	{Name: "NoncurrentDays", JSONPath: "NoncurrentDays"},
	{Name: "AbortDays", JSONPath: "AbortDays"},
}

type ruleInfo struct {
	ID                        string `json:"ID"`
	Prefix                    string `json:"Prefix"`
	Status                    string `json:"Status"`
	ExpirationDays            string `json:"ExpirationDays"`
	ExpirationDate            string `json:"ExpirationDate"`
	ExpiredObjectDeleteMarker string `json:"ExpiredObjectDeleteMarker"`
	NoncurrentDays            string `json:"NoncurrentDays"`
	AbortDays                 string `json:"AbortDays"`
}

func int32PtrToStr(v *int32) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}

func boolPtrToStr(v *bool) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%t", *v)
}

func LifecycleCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lifecycle",
			Aliases:          []string{"lc"},
			Short:            "Bucket lifecycle operations for contract-owned object storage",
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
