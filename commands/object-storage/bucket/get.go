package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type bucketInfo struct {
	Name         string    `json:"Name"`
	CreationDate time.Time `json:"CreationDate"`
	Region       string    `json:"Region"`
}

// getBucketInfo fetches metadata for a single bucket by name.
// S3 has no single-bucket GET endpoint; ListBuckets is the only way to retrieve creation date.
func getBucketInfo(ctx context.Context, name string) (*bucketInfo, error) {
	result, _, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.ListBuckets(ctx).Execute()
	if err != nil {
		return nil, err
	}

	var found *bucketInfo
	for _, b := range result.GetBuckets() {
		if b.GetName() == name {
			found = &bucketInfo{
				Name:         b.GetName(),
				CreationDate: b.GetCreationDate(),
			}
			break
		}
	}
	if found == nil {
		return nil, fmt.Errorf("bucket %q not found", name)
	}

	loc, _, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.GetBucketLocation(ctx, name).Execute()
	if err == nil && loc != nil {
		found.Region = loc.GetLocationConstraint()
	}

	return found, nil
}

func GetBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get details of a contract-owned bucket",
		Example:   "ionosctl object-storage bucket get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			found, err := getBucketInfo(c.Context, name)
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, found, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to retrieve", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
