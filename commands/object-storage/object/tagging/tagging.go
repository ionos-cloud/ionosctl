package objecttagging

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagKey       = "key"
	flagKeyShort  = "k"
	flagVersionId = "version-id"
)

var allCols = []table.Column{
	{Name: "Key", JSONPath: "Key", Default: true},
	{Name: "Value", JSONPath: "Value", Default: true},
}

type tagInfo struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func ObjectTaggingCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "tagging",
			Aliases: []string{"tag"},
			Short:   "Manage the tag set (key/value tags) on an object",
			Long: `Manage the tag set attached to an object.

Object tags are up to 10 key/value string pairs stored alongside an object. Unlike metadata (which is fixed at upload time), tags can be added, changed or removed at any point without re-uploading the object. They are commonly used for cost allocation, for targeting objects in lifecycle rules, and as conditions in access policies.

The tag set is managed as a whole: "put" REPLACES the entire set (it is not a merge), "get" reads it, and "delete" removes all tags at once. On versioning-enabled buckets, tags belong to a specific object version; pass --version-id to act on a version other than the current one.`,
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
