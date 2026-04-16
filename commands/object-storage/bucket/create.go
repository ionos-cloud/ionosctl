package bucket

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CreateBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a contract-owned bucket",
		Example:   "ionosctl object-storage bucket create --name my-bucket\nionosctl object-storage bucket create --name my-bucket --region eu-central-3",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			location := viper.GetString(constants.FlagLocation)

			cfg := objectstorage.NewCreateBucketConfiguration()
			cfg.SetLocationConstraint(location)

			_, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.CreateBucket(c.Context, name).
				CreateBucketConfiguration(*cfg).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q created successfully in region %q\n", name, location)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to create", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
