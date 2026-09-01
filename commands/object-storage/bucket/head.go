package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var headCols = []table.Column{
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "Status", JSONPath: "Status", Default: true},
	{Name: "Region", JSONPath: "Region", Default: true},
}

type headResult struct {
	Name   string `json:"Name"`
	Status string `json:"Status"`
	Region string `json:"Region"`
}

func HeadBucketCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "head",
		Aliases:   []string{"hd"},
		ShortDesc: "Check if a bucket exists and you have access",
		LongDesc: `Perform an S3 HeadBucket check: a lightweight probe that confirms the bucket exists AND that your credentials are allowed to access it, without transferring any bucket contents.

The distinction from 'bucket get' matters: HeadBucket returns HTTP status only. A 404 means the bucket does not exist; a 403 means it exists but your credentials lack access; success means it exists and is accessible. On success this also resolves and reports the bucket's real region rather than echoing back --location.`,
		Example: "ionosctl object-storage bucket head --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			_, err := client.MustObjectStorage().ObjectStorageClient.BucketsApi.HeadBucket(c.Context, name).Execute()
			if err != nil {
				return fmt.Errorf("checking if bucket %q exists: %w", name, err)
			}

			// Report the bucket's actual region rather than echoing --location,
			// which may be unset or differ from where the bucket really lives.
			region := ""
			if loc, _, locErr := client.MustObjectStorage().ObjectStorageClient.BucketsApi.
				GetBucketLocation(c.Context, name).Execute(); locErr == nil {
				region = loc.GetLocationConstraint()
			}

			result := headResult{
				Name:   name,
				Status: "exists and is accessible",
				Region: region,
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(headCols, result, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to check for existence and access", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, table.ColsMessage(headCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(headCols), cobra.ShellCompDirectiveNoFileComp
		})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
