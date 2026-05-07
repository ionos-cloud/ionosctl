package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

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
		Example:   "ionosctl object-storage bucket list\nionosctl object-storage bucket list --location eu-central-3",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			result, _, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.ListBuckets(c.Context).Execute()
			if err != nil {
				return err
			}

			filterByLocation := viper.IsSet(constants.FlagLocation)
			filterRegion := viper.GetString(constants.FlagLocation)

			var buckets []bucketInfo
			for _, b := range result.GetBuckets() {

				bi := bucketInfo{
					Name:         b.GetName(),
					CreationDate: b.GetCreationDate(),
				}

				loc, apiResp, locErr := client.MustObjectStorage().ObjectStorageClient.BucketsApi.GetBucketLocation(c.Context, b.GetName()).Execute()
				if locErr != nil {
					// When filtering by location, a failed lookup would silently
					// exclude the bucket. Fail loudly so the user knows results
					// may be incomplete.
					if filterByLocation {
						return fmt.Errorf("cannot filter by location: failed to get location for bucket %q: %w", b.GetName(), locErr)
					}

					c.Verbose("failed to get location for bucket %q: %v", b.GetName(), locErr)
					if apiResp != nil {
						bi.Region = fmt.Sprintf("<%s>", apiResp.Status)
					} else {
						bi.Region = "<unknown>"
					}
				} else {
					bi.Region = loc.GetLocationConstraint()
				}

				if filterByLocation && bi.Region != filterRegion {
					continue
				}

				buckets = append(buckets, bi)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, buckets, cols))
		},
		InitClient: false,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
