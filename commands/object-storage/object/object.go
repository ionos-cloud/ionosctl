package object

import (
	"github.com/spf13/cobra"

	legalhold "github.com/ionos-cloud/ionosctl/v6/commands/object-storage/object/legal-hold"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/object/retention"
	objecttagging "github.com/ionos-cloud/ionosctl/v6/commands/object-storage/object/tagging"
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
			Use:     "object",
			Aliases: []string{"obj"},
			Short:   "Manage objects (files) stored in IONOS Object Storage buckets",
			Long: `Manage the objects (files) stored inside IONOS Object Storage buckets.

IONOS Object Storage is S3-compatible; these commands map onto the S3 object API, so S3 semantics apply.

Data model:
  - A bucket is a flat namespace of objects. An OBJECT is identified by its KEY (a UTF-8 string, e.g. "photos/2025/image.jpg"). The "/" in a key is only a naming convention: there are no real directories. Use --prefix on "list" to browse a pseudo-folder.
  - Every object carries a value (its bytes), a content-type, an ETag (a checksum/opaque identity token), a size and a last-modified time. "head" returns this metadata without downloading the bytes; "get" downloads the bytes.
  - VERSIONING: when a bucket has versioning enabled, overwriting or deleting a key does not destroy older data. Each write produces a distinct VERSION ID; a delete inserts a "delete marker" that hides the key while leaving prior versions intact. Most subcommands accept --version-id to target one specific version instead of the current one.

Sub-resources:
  - tagging     key/value tags on an object, used for cost allocation, lifecycle rules and access policies.
  - retention   Object Lock WORM retention: protect an object from deletion until a date (see "retention").
  - legal-hold  Object Lock legal hold: an on/off protection independent of any retain-until date (see "legal-hold").

Retention vs legal hold: retention protects an object until a specific date and then automatically lapses; a legal hold has no expiry and protects the object until it is explicitly turned OFF. They are independent - an object can have both, either, or neither, and it stays locked while ANY protection is in force. Both require the bucket to have been created with Object Lock enabled (Object Lock cannot be added to an existing bucket, and it forces versioning on).`,
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

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(PutCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(DeleteCmd())
	cmd.AddCommand(HeadCmd())
	cmd.AddCommand(CopyCmd())
	cmd.AddCommand(retention.Root())
	cmd.AddCommand(legalhold.Root())
	cmd.AddCommand(objecttagging.ObjectTaggingCmd())

	return cmd
}
