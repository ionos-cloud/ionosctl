package objectstorage

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/cors"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/encryption"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/object"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/policy"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/tagging"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "object-storage",
			Aliases:          []string{"os"},
			Short:            "Object Storage operations for contract-owned buckets",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(bucket.BucketCommand())
	cmd.AddCommand(cors.CorsCmd())
	cmd.AddCommand(encryption.EncryptionCmd())
	cmd.AddCommand(object.ObjectCommand())
	cmd.AddCommand(policy.PolicyCmd())
	cmd.AddCommand(tagging.TaggingCmd())
	return cmd
}
