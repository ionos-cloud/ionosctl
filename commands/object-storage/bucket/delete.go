package bucket

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/viper"
)

func DeleteBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a contract-owned bucket",
		LongDesc:  "Delete a contract-owned bucket. The bucket must be empty before it can be deleted.",
		Example:   "ionosctl object-storage bucket delete --name my-bucket\nionosctl object-storage bucket delete --name my-bucket --region eu-central-3",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			region := viper.GetString(core.GetFlagName(c.NS, constants.FlagS3Region))

			if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete bucket %q", name), viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			s3, err := client.GetObjectStorageClient(region)
			if err != nil {
				return err
			}

			_, err = s3.BucketsApi.DeleteBucket(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q deleted successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to delete", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagS3Region, "r", "eu-central-3", "Region of the bucket (e.g. eu-central-3)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
