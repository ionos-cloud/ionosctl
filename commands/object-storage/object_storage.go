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

	cmd = core.WithRegionalConfigOverride(cmd,
		[]string{fileconfiguration.ObjectStorage},
		constants.ObjectStorageApiRegionalURL,
		constants.ObjectStorageLocations,
	)

	// Document the default. The shared flag has no cobra default (list queries all
	// locations when unset), but single-resource commands fall back to the first
	// location, so note it here rather than as a misleading global cobra default.
	if f := cmd.Command.PersistentFlags().Lookup(constants.FlagLocation); f != nil {
		f.Usage += ". Defaults to " + constants.ObjectStorageLocations[0]
	}

	return cmd
}
