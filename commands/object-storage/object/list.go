package object

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

const (
	flagPrefix  = "prefix"
	flagMaxKeys = "max-keys"
)

var listCols = []table.Column{
	{Name: "Key", JSONPath: "Key", Default: true},
	{Name: "Size", JSONPath: "Size", Default: true},
	{Name: "LastModified", JSONPath: "LastModified", Default: true},
	{Name: "StorageClass", JSONPath: "StorageClass", Default: true},
	{Name: "ETag", JSONPath: "ETag"},
}

type listObjectInfo struct {
	Key          string    `json:"Key"`
	Size         int32     `json:"Size"`
	LastModified time.Time `json:"LastModified"`
	StorageClass string    `json:"StorageClass"`
	ETag         string    `json:"ETag"`
}

func ListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List objects in a bucket",
		Example:   "ionosctl object-storage object list --name my-bucket\nionosctl object-storage object list --name my-bucket --prefix photos/ --max-keys 100",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			prefix := viper.GetString(core.GetFlagName(c.NS, flagPrefix))
			maxKeys := viper.GetInt32(core.GetFlagName(c.NS, flagMaxKeys))

			s3Regional, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			var allObjects []listObjectInfo
			var continuationToken string
			noLimit := maxKeys <= 0
			remaining := maxKeys

			for {
				req := s3Regional.ObjectsApi.ListObjectsV2(c.Context, name)

				if prefix != "" {
					req = req.Prefix(prefix)
				}

				if !noLimit {
					req = req.MaxKeys(min(remaining, 1000))
				}

				if continuationToken != "" {
					req = req.ContinuationToken(continuationToken)
				}

				result, _, err := req.Execute()
				if err != nil {
					return err
				}

				for _, obj := range result.Contents {
					info := listObjectInfo{
						Key:  obj.GetKey(),
						Size: obj.GetSize(),
						ETag: obj.GetETag(),
					}
					if obj.LastModified != nil {
						info.LastModified = obj.LastModified.Time
					}
					if obj.StorageClass != nil {
						info.StorageClass = string(*obj.StorageClass)
					}
					allObjects = append(allObjects, info)
				}

				if !noLimit {
					remaining -= int32(len(result.Contents))
					if remaining <= 0 {
						break
					}
				}

				if !result.IsTruncated {
					break
				}

				if result.NextContinuationToken != nil {
					continuationToken = *result.NextContinuationToken
				} else {
					break
				}
			}

			if len(allObjects) == 0 {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), "No objects found")
				return nil
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(listCols, allObjects, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagPrefix, "p", "", "Filter objects by key prefix (e.g. photos/)")
	cmd.AddInt32Flag(flagMaxKeys, "", 1000, "Maximum number of objects to return (0 for no limit)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
