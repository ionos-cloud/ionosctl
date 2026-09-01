package object

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type headObjectInfo struct {
	Key           string `json:"Key"`
	ContentType   string `json:"ContentType"`
	ContentLength string `json:"ContentLength"`
	LastModified  string `json:"LastModified"`
	ETag          string `json:"ETag"`
}

func HeadCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "head",
		Aliases:   []string{"hd"},
		ShortDesc: "Get an object's metadata without downloading its bytes",
		LongDesc: `Fetch an object's metadata without transferring its contents (an S3 HEAD request).

Returns the key, content-type, content-length (size in bytes), last-modified time and ETag. This is the cheap way to check whether an object exists, how large it is, or what type it is - use "object get" when you actually need the bytes.`,
		Example: "ionosctl object-storage object head --name my-bucket --key photos/image.jpg",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))

			_, apiResp, err := client.MustObjectStorage().ObjectStorageClient.ObjectsApi.HeadObject(c.Context, name, key).Execute()
			if err != nil {
				return err
			}

			info := headObjectInfo{
				Key:           key,
				ContentType:   apiResp.Header.Get("Content-Type"),
				ContentLength: apiResp.Header.Get("Content-Length"),
				LastModified:  apiResp.Header.Get("Last-Modified"),
				ETag:          apiResp.Header.Get("ETag"),
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(headCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key (full name) of the object to inspect", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
