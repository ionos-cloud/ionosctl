package objectstorage

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/bucket"
	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/object"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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
	cmd.AddCommand(object.ObjectCommand())

	return core.WithRegionalConfigOverride(cmd,
		[]string{fileconfiguration.ObjectStorage},
		constants.ObjectStorageApiRegionalURL,
		constants.ObjectStorageLocations,
	)
}
