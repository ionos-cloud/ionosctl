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
		Example:   "ionosctl object-storage bucket create --name my-bucket\nionosctl object-storage bucket create --name my-bucket --location eu-central-3\nionosctl object-storage bucket create --name my-bucket --object-lock",
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to create", core.RequiredFlagOption())
	cmd.AddBoolFlag(flagObjectLock, "", false, "Enable Object Lock on the new bucket (cannot be changed after creation)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
