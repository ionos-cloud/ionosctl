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
			Use:     "lifecycle",
			Aliases: []string{"lc"},
			Short:   "Manage automatic object expiration/cleanup rules for a bucket",
			Long: `Manage a bucket's lifecycle configuration: rules that let the service automatically expire and clean up objects over time, so you do not have to delete them by hand. Each rule targets objects by key Prefix (empty prefix = the whole bucket) and can:
  - Expiration: delete current objects a set number of Days after creation (or on a fixed Date).
  - NoncurrentVersionExpiration: on a versioned bucket, delete OLD (noncurrent) versions a number of days after they stop being current.
  - AbortIncompleteMultipartUpload: discard the parts of multipart uploads that were never completed after N days (reclaims storage you are still billed for).
  - ExpiredObjectDeleteMarker: clean up dangling delete markers left behind on versioned buckets.

Versioning interplay: on a versioned bucket, plain Expiration does NOT free storage - it just adds a delete marker over the current version while all prior versions remain. To actually reclaim space you must also use NoncurrentVersionExpiration. Rules run asynchronously (typically once per day), so deletion is not instantaneous.

The bucket holds one lifecycle configuration; 'put' replaces it wholesale, 'get' lists the current rules, 'delete' removes lifecycle management entirely.`,
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
