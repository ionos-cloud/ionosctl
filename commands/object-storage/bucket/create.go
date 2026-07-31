package bucket

import (
	"context"

	"github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const flagObjectLock = "object-lock"

func CreateBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a contract-owned bucket",
		LongDesc: `Create a new S3-compatible bucket owned by your contract.

The bucket name must be globally unique across the whole Object Storage service (not just your account) and follow S3 bucket-naming rules (3-63 chars, lowercase letters, numbers, dots and hyphens; DNS-compatible). The bucket is created in a single location and stays there for its lifetime.

--location selects the region the bucket is created in and sets its Locationconstraint; it also decides which regional endpoint the request is signed for, so the two always stay in sync. When --location is omitted, the first Object Storage location (` + "`" + `` + constants.ObjectStorageLocations[0] + `` + "`" + `) is used.

--object-lock enables WORM (Write-Once-Read-Many) Object Lock on the bucket. This is a create-time-only decision: it CANNOT be turned on (or off) after the bucket exists. Enabling it also implicitly enables versioning. After creation, define the default retention with ` + "`bucket object-lock put`" + `.`,
		Example: `# Create a bucket in the default location
ionosctl object-storage bucket create --name my-bucket

# Create a bucket in a specific region
ionosctl object-storage bucket create --name my-bucket --location eu-central-3

# Create a bucket with WORM Object Lock enabled (irreversible; implies versioning)
ionosctl object-storage bucket create --name my-bucket --object-lock`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			// --location is optional here: when unset it falls back to the first
			// Object Storage location, matching the endpoint that PersistentPreRunE
			// resolves for the same unset case. Sourced explicitly rather than from
			// the flag default so the LocationConstraint and endpoint stay in sync.
			location := viper.GetString(constants.FlagLocation)
			if location == "" {
				location = constants.ObjectStorageLocations[0]
			}
			objectLock := viper.GetBool(core.GetFlagName(c.NS, flagObjectLock))

			cfg := objectstorage.NewCreateBucketConfiguration()
			cfg.SetLocationConstraint(location)

			req := client.MustObjectStorage().ObjectStorageClient.BucketsApi.CreateBucket(c.Context, name).
				CreateBucketConfiguration(*cfg).XAmzBucketObjectLockEnabled(objectLock)
			_, err := req.Execute()
			if err != nil {
				return err
			}

			info, err := getBucketInfo(c.Context, name)
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Globally-unique bucket name (3-63 chars, lowercase, DNS-compatible). Must not already exist anywhere in the service", core.RequiredFlagOption())
	cmd.AddBoolFlag(flagObjectLock, "", false, "Enable WORM Object Lock at creation. Irreversible and implicitly enables versioning; set the default retention afterwards with 'bucket object-lock put'. Default: false")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
