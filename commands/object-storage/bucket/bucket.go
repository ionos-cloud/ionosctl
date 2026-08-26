package bucket

import (
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/cors"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/encryption"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/lifecycle"
	objectlock "github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/object-lock"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/policy"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/publicaccessblock"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/tagging"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket/versioning"
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
			Use:     "bucket",
			Aliases: []string{"b"},
			Short:   "Create, inspect, and configure S3-compatible buckets",
			Long: `Manage contract-owned IONOS Object Storage buckets. A bucket is a container for objects, uniquely named across the whole service and permanently bound to the location it was created in. This subtree also configures the per-bucket features layered on top of the bucket, each addressed by bucket name:

  create/get/head/list/delete   Lifecycle of the bucket itself (head only checks existence + access; get/list resolve the bucket's region).
  versioning                    Keep multiple versions of each object (Enabled/Suspended). Once enabled it can only be suspended, never turned fully off.
  object-lock                   Write-Once-Read-Many (WORM) retention. Requires the bucket to have been created with --object-lock and requires versioning.
  lifecycle                     Rules that expire/clean up objects (and noncurrent versions / aborted multipart uploads) automatically over time.
  encryption                    Default server-side encryption (SSE) applied to new objects when the request does not specify its own.
  policy                        JSON access policy (S3/IAM syntax) attached to the bucket, plus a status check for whether it makes the bucket public.
  public-access-block           Account/bucket guardrails that override ACLs and policies to keep a bucket private.
  cors                          Cross-Origin Resource Sharing rules governing browser access from other origins.
  tagging                       Key/value tags on the bucket (for cost allocation / organization).

These features interact (e.g. object-lock requires versioning, and public-access-block overrides ACLs and policies); see each subcommand's help for the details.`,
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
	cmd.AddCommand(DeleteBucketCmd())
	cmd.AddCommand(versioning.Root())
	cmd.AddCommand(objectlock.Root())
	cmd.AddCommand(cors.CorsCmd())
	cmd.AddCommand(encryption.EncryptionCmd())
	cmd.AddCommand(tagging.TaggingCmd())
	cmd.AddCommand(policy.PolicyCmd())
	cmd.AddCommand(lifecycle.LifecycleCmd())
	cmd.AddCommand(publicaccessblock.PublicAccessBlockCmd())
	return cmd
}
