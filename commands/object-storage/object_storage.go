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
			Use:     "object-storage",
			Aliases: []string{"os"},
			Short:   "Manage IONOS Object Storage (S3-compatible) buckets and objects",
			Long: `IONOS Object Storage is an S3-compatible object store. Data is organized into buckets (flat namespaces that hold objects) served from regional, per-location endpoints of the form https://s3.<location>.ionoscloud.com. Because the API mirrors the AWS S3 API, standard S3 semantics, request/response shapes and tooling apply.

Authentication differs from the rest of ionosctl. Most IONOS APIs use your IONOS token, but Object Storage is reached with S3 credentials: an access key and secret key (an S3 SigV4 signature), scoped to a region. Generate these keys in the DCD or via the object-storage-management API, then configure them for ionosctl (see the config/credentials setup). A plain IONOS token alone will NOT authenticate S3 calls.

Region/location: every request is signed for and routed to one location. Commands accept --location to pick the region (` + "`eu-central-3`" + `, ` + "`eu-central-4`" + `, ` + "`us-central-1`" + `); when omitted, single-resource commands default to ` + "`" + `` + constants.ObjectStorageLocations[0] + `` + "`" + ` while ` + "`bucket list`" + ` queries all locations. Bucket names are globally unique across the service and a bucket permanently lives in the location it was created in.

This command tree covers "contract-owned" buckets (billed to and owned by your IONOS contract). Use the ` + "`bucket`" + ` subtree to manage buckets and their per-bucket configurations, and the ` + "`object`" + ` subtree to manage the objects inside them.`,
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
