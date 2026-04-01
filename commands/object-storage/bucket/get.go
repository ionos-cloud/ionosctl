package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
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

			s3, err := client.GetObjectStorageClient("")
			if err != nil {
				return err
			}

			// S3 has no API to get a single bucket's metadata (creation date).
			// ListBuckets is the only way to retrieve it.
			result, _, err := s3.BucketsApi.ListBuckets(context.Background()).Execute()
			if err != nil {
				return err
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
				return fmt.Errorf("bucket %q not found", name)
			}

			loc, _, err := s3.BucketsApi.GetBucketLocation(context.Background(), name).Execute()
			if err == nil && loc != nil {
				found.Region = loc.GetLocationConstraint()
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, found, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to retrieve", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
