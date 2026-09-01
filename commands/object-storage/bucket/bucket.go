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
			Long: `Manage contract-owned IONOS Object Storage buckets. A bucket is a container for objects, uniquely named across the whole service and permanently bound to the location it was created in. This subtree also configures the per-bucket features layered on top of the bucket (versioning, object-lock, lifecycle, encryption, policy, public-access-block, CORS and tagging), each addressed by bucket name.

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
