package object

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func GetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Download an object's bytes to a local file",
		LongDesc: `Download the bytes of an object to a local file.

By default the file is written to the current directory under the basename of the key (the part after the last "/", so photos/image.jpg is saved as image.jpg). Use --destination to write elsewhere or under a different name.

On a versioning-enabled bucket, the current (latest) version is downloaded unless you pass --version-id to fetch a specific historical version. Version IDs are shown by "object list" (with the appropriate columns) and in the S3 version listing.

To fetch only an object's metadata (size, content-type, ETag, last-modified) without downloading the bytes, use "object head" instead.`,
		Example: `# Download to ./image.jpg (basename of the key) in the current directory
ionosctl object-storage object get --name my-bucket --key photos/image.jpg

# Download to an explicit local path
ionosctl object-storage object get --name my-bucket --key photos/image.jpg --destination ./local-image.jpg

# Download a specific historical version
ionosctl object-storage object get --name my-bucket --key photos/image.jpg --version-id <version-id>`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			destination := viper.GetString(core.GetFlagName(c.NS, flagDestination))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			if destination == "" {
				destination = filepath.Base(key)
			}

			req := client.MustObjectStorage().ObjectStorageClient.ObjectsApi.GetObject(c.Context, name, key)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			tmpFile, _, err := req.Execute()
			if err != nil {
				return err
			}
			defer func() {
				tmpFile.Close()
				os.Remove(tmpFile.Name())
			}()

			outFile, err := os.Create(destination)
			if err != nil {
				return fmt.Errorf("creating destination file %q: %w", destination, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tmpFile); err != nil {
				return fmt.Errorf("writing to %q: %w", destination, err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Object %q downloaded to %q\n", key, destination)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key (full name) of the object to download", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagDestination, "d", "", "Local file path to write the download to. Defaults to the basename of the key (the part after the last \"/\") in the current directory")
	cmd.AddStringFlag(flagVersionId, "", "", "Download this specific object version instead of the current one (versioned buckets only). Defaults to the latest version")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
