package objectstorage

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/policy"
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
	cmd.AddCommand(policy.PolicyCmd())
	return cmd
}
