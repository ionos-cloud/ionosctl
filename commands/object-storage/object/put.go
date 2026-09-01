package object

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Upload a local file as an object",
		LongDesc: `Upload a local file into a bucket, stored under the key you choose with --key.

The key is the object's full name in the bucket, including any "/" separators you want to appear as a pseudo-folder path (e.g. photos/2025/image.jpg). The key does NOT have to match the source filename. If the key already exists it is overwritten; on a versioning-enabled bucket the overwrite creates a new version rather than destroying the old bytes.

--content-type sets the object's MIME type, which the service stores and returns as the Content-Type header on later downloads (it drives how browsers and clients interpret the bytes). If omitted, it is auto-detected from the source file extension, falling back to application/octet-stream when the extension is unknown.

The whole file is read from --source and uploaded as the object body.`,
		Example: `# Upload a file, letting the content-type be auto-detected from the extension
ionosctl object-storage object put --name my-bucket --key photos/image.jpg --source ./image.jpg

# Upload under a different key and set an explicit content-type
ionosctl object-storage object put --name my-bucket --key exports/report.json --source ./out.dat --content-type application/json`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, flagSource)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			source := viper.GetString(core.GetFlagName(c.NS, flagSource))
			contentType := viper.GetString(core.GetFlagName(c.NS, flagContentType))

			file, err := os.Open(source)
			if err != nil {
				return fmt.Errorf("opening source file %q: %w", source, err)
			}
			defer file.Close()

			if contentType == "" {
				contentType = mime.TypeByExtension(filepath.Ext(source))
				if contentType == "" {
					contentType = "application/octet-stream"
				}
			}

			_, err = client.MustObjectStorage().ObjectStorageClient.ObjectsApi.PutObject(c.Context, name, key).
				Body(file).
				ContentType(contentType).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Object %q uploaded to bucket %q\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the destination bucket to upload into", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key to store the object under, i.e. its full name in the bucket. Use \"/\" to build a pseudo-folder path (e.g. photos/image.jpg). An existing key is overwritten (creating a new version on versioned buckets)", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagSource, flagSourceShort, "", "Path to the local file whose bytes are uploaded as the object body", core.RequiredFlagOption())
	cmd.AddStringFlag(flagContentType, "", "", "MIME type stored with the object and returned as its Content-Type on download (e.g. image/jpeg). Auto-detected from the --source file extension when omitted, defaulting to application/octet-stream if unknown")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
