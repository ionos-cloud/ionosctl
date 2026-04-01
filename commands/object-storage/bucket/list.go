package bucket

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func ListBucketsCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List all contract-owned buckets",
		Example:   "ionosctl object-storage bucket list\nionosctl object-storage bucket list --region eu-central-3",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			s3, err := client.GetObjectStorageClient("")
			if err != nil {
				return err
			}

			result, _, err := s3.BucketsApi.ListBuckets(context.Background()).Execute()
			if err != nil {
				return err
			}

			filterRegion := viper.GetString(core.GetFlagName(c.NS, constants.FlagS3Region))

			var buckets []bucketInfo
			for _, b := range result.GetBuckets() {
				bi := bucketInfo{
					Name:         b.GetName(),
					CreationDate: b.GetCreationDate(),
				}

				loc, _, locErr := s3.BucketsApi.GetBucketLocation(context.Background(), b.GetName()).Execute()
				if locErr == nil && loc != nil {
					bi.Region = loc.GetLocationConstraint()
				}

				if filterRegion != "" && bi.Region != filterRegion {
					continue
				}

				buckets = append(buckets, bi)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, buckets, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagS3Region, "r", "", "Filter buckets by region (e.g. eu-central-3)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagS3Region, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.Regions(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
