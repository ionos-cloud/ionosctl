package object

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type copyObjectInfo struct {
	ETag         string `json:"ETag"`
	LastModified string `json:"LastModified"`
}

func CopyCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "copy",
		Aliases:   []string{"cp"},
		ShortDesc: "Server-side copy an object to a new key or bucket",
		LongDesc: `Copy an object to a new key, within the same bucket or into another bucket.

The copy happens entirely SERVER-SIDE: the bytes are never downloaded to your machine and re-uploaded. Both --name (the destination bucket) and --copy-source (the source) must live in the same region/endpoint.

--copy-source names the object being copied and must be in the form /source-bucket/source-key (a leading slash, then the bucket, then the key). The destination is --name plus --key. Copying onto an existing key overwrites it; on a versioning-enabled destination bucket the copy becomes a new version.

Common uses: rename or move an object (copy to the new key, then delete the old one), or duplicate an object into another bucket. The command prints the new object's ETag and LastModified on success.`,
		Example: `# Copy within the same bucket (e.g. to "rename" a key - delete the old one afterwards)
ionosctl object-storage object copy --name my-bucket --key photos/renamed.jpg --copy-source /my-bucket/photos/image.jpg

# Copy across buckets
ionosctl object-storage object copy --name backup-bucket --key photos/image.jpg --copy-source /my-bucket/photos/image.jpg`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, flagCopySource)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			copySource := viper.GetString(core.GetFlagName(c.NS, flagCopySource))

			result, _, err := client.MustObjectStorage().ObjectStorageClient.ObjectsApi.CopyObject(c.Context, name, key).
				XAmzCopySource(copySource).
				Execute()
			if err != nil {
				return err
			}

			info := copyObjectInfo{
				ETag: result.GetETag(),
			}
			if lm := result.GetLastModified(); !lm.IsZero() {
				info.LastModified = lm.Format(time.RFC3339)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(copyCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Destination bucket to copy into (must be in the same region as the source)", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key to store the copy under in the destination bucket (overwrites if it already exists)", core.RequiredFlagOption())
	cmd.AddStringFlag(flagCopySource, "", "", "The source object to copy, in the form /source-bucket/source-key (leading slash required)", core.RequiredFlagOption())

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, table.ColsMessage(copyCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(copyCols), cobra.ShellCompDirectiveNoFileComp
		})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
