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
			Use:              "bucket",
			Aliases:          []string{"b"},
			Short:            "Bucket operations for contract-owned object storage",
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
