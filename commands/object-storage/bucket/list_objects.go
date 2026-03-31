package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const (
	flagPrefix  = "prefix"
	flagMaxKeys = "max-keys"
)

var objectCols = []table.Column{
	{Name: "Key", JSONPath: "Key", Default: true},
	{Name: "Size", JSONPath: "Size", Default: true},
	{Name: "LastModified", JSONPath: "LastModified", Default: true},
	{Name: "StorageClass", JSONPath: "StorageClass", Default: true},
	{Name: "ETag", JSONPath: "ETag"},
}

type objectInfo struct {
	Key          string    `json:"Key"`
	Size         int32     `json:"Size"`
	LastModified time.Time `json:"LastModified"`
	StorageClass string    `json:"StorageClass"`
	ETag         string    `json:"ETag"`
}

func ListObjectsCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket",
		Verb:      "list-objects",
		Aliases:   []string{"lo"},
		ShortDesc: "List objects in a bucket",
		Example:   "ionosctl object-storage bucket list-objects --name my-bucket\nionosctl object-storage bucket list-objects --name my-bucket --prefix photos/ --max-keys 100",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			prefix := viper.GetString(core.GetFlagName(c.NS, flagPrefix))
			maxKeys := viper.GetInt32(core.GetFlagName(c.NS, flagMaxKeys))

			s3, err := client.GetObjectStorageClient("")
			if err != nil {
				return err
			}

			// Resolve the bucket's region to avoid redirect loops.
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

			var allObjects []objectInfo
			var continuationToken string
			remaining := maxKeys

			for {
				req := s3Regional.ObjectsApi.ListObjectsV2(context.Background(), name)

				if prefix != "" {
					req = req.Prefix(prefix)
				}

				pageSize := remaining
				if pageSize > 1000 {
					pageSize = 1000
				}
				req = req.MaxKeys(pageSize)

				if continuationToken != "" {
					req = req.ContinuationToken(continuationToken)
				}

				result, _, err := req.Execute()
				if err != nil {
					return err
				}

				for _, obj := range result.Contents {
					info := objectInfo{
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

				remaining -= int32(len(result.Contents))

				if !result.IsTruncated || remaining <= 0 {
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
			return c.Out(table.Sprint(objectCols, allObjects, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	cmd.AddStringFlag(flagPrefix, "p", "", "Filter objects by key prefix (e.g. photos/)")
	cmd.AddInt32Flag(flagMaxKeys, "", 1000, "Maximum number of objects to return")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
