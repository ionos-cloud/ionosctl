package object

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagKey         = "key"
	flagKeyShort    = "k"
	flagSource      = "source"
	flagSourceShort = "s"
	flagDestination = "destination"
	flagVersionId   = "version-id"
	flagCopySource  = "copy-source"
	flagContentType = "content-type"
)

var headCols = []table.Column{
	{Name: "Key", JSONPath: "Key", Default: true},
	{Name: "ContentType", JSONPath: "ContentType", Default: true},
	{Name: "ContentLength", JSONPath: "ContentLength", Default: true},
	{Name: "LastModified", JSONPath: "LastModified", Default: true},
	{Name: "ETag", JSONPath: "ETag", Default: true},
}

var copyCols = []table.Column{
	{Name: "ETag", JSONPath: "ETag", Default: true},
	{Name: "LastModified", JSONPath: "LastModified", Default: true},
}

func ObjectCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "object",
			Aliases:          []string{"obj"},
			Short:            "Object operations for contract-owned object storage",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(headCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(headCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(PutCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(HeadCmd())
	cmd.AddCommand(CopyCmd())

	return cmd
}
