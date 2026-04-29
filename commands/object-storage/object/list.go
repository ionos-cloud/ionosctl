package object

import (
	"context"
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
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
	Size         string    `json:"Size"`
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
		CmdRun:     runListObjects,
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagPrefix, "p", "", "Filter objects by key prefix (e.g. photos/)")
	cmd.AddInt32Flag(flagMaxKeys, "", 1000, "Maximum number of objects to return (0 for no limit)")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, table.ColsMessage(listCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(listCols), cobra.ShellCompDirectiveNoFileComp
		})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

func runListObjects(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	prefix := viper.GetString(core.GetFlagName(c.NS, flagPrefix))
	maxKeys := viper.GetInt32(core.GetFlagName(c.NS, flagMaxKeys))

	s3 := client.MustObjectStorage().ObjectStorageClient

	var allObjects []listObjectInfo
	var continuationToken string
	noLimit := maxKeys <= 0
	remaining := maxKeys

	for {
		req := s3.ObjectsApi.ListObjectsV2(c.Context, name)

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

		allObjects = append(allObjects, convertObjects(result.Contents)...)

		if !noLimit {
			remaining -= int32(len(result.Contents))
			if remaining <= 0 {
				break
			}
		}

		if !result.IsTruncated || result.NextContinuationToken == nil {
			break
		}
		continuationToken = *result.NextContinuationToken
	}

	if len(allObjects) == 0 {
		fmt.Fprintln(c.Command.Command.OutOrStdout(), "No objects found")
		return nil
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(listCols, allObjects, cols))
}

func convertObjects(objects []objectstorage.Object) []listObjectInfo {
	infos := make([]listObjectInfo, 0, len(objects))
	for _, obj := range objects {
		info := listObjectInfo{
			Key:  obj.GetKey(),
			Size: humanize.IBytes(uint64(int64(obj.GetSize()))),
			ETag: obj.GetETag(),
		}
		if obj.LastModified != nil {
			info.LastModified = obj.LastModified.Time
		}
		if obj.StorageClass != nil {
			info.StorageClass = string(*obj.StorageClass)
		}
		infos = append(infos, info)
	}
	return infos
}
