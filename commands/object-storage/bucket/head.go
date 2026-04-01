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
		Example:   "ionosctl object-storage bucket head --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3Regional, region, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			_, err = s3Regional.BucketsApi.HeadBucket(context.Background(), name).Execute()
			if err != nil {
				return err
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket to check", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
