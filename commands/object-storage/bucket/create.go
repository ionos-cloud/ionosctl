package bucket

import (
	"context"
	"fmt"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
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
			region := viper.GetString(core.GetFlagName(c.NS, constants.FlagS3Region))

			s3, err := client.GetObjectStorageClient(region)
			if err != nil {
				return err
			}

			cfg := objectstorage.NewCreateBucketConfiguration()
			cfg.SetLocationConstraint(region)

			_, err = s3.BucketsApi.CreateBucket(c.Context, name).
				CreateBucketConfiguration(*cfg).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q created successfully in region %q\n", name, region)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to create", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagS3Region, "r", "eu-central-3", "Region where the bucket will be created (e.g. eu-central-3)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagS3Region, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.Regions(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
