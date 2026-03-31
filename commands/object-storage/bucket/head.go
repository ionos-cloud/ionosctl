package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func HeadBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "head",
		Aliases:   []string{"h"},
		ShortDesc: "Check if a bucket exists and you have access",
		Example:   "ionosctl object-storage bucket head --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			// Use a global endpoint to discover the bucket's actual region,
			// then issue HeadBucket against that region's endpoint to avoid
			// redirect loops caused by cross-region requests.
			s3, err := client.GetObjectStorageClient("")
			if err != nil {
				return err
			}

			loc, _, err := s3.BucketsApi.GetBucketLocation(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			region := ""
			if loc != nil {
				region = loc.GetLocationConstraint()
			}

			s3Regional, err := client.GetObjectStorageClient(region)
			if err != nil {
				return err
			}

			_, err = s3Regional.BucketsApi.HeadBucket(context.Background(), name).Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket %q exists and is accessible\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to check", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
