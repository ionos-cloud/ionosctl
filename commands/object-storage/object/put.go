package object

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
		ShortDesc: "Upload a file as an object",
		Example:   "ionosctl object-storage object put --name my-bucket --key photos/image.jpg --source ./image.jpg",
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

			s3Regional, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			_, err = s3Regional.ObjectsApi.PutObject(c.Context, name, key).
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key (path in the bucket)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagKey, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		bucket := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName))
		return completer.ObjectKeys(bucket), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagSource, flagSourceShort, "", "Path to the local file to upload", core.RequiredFlagOption())
	cmd.AddStringFlag(flagContentType, "", "", "MIME type of the object (auto-detected from file extension if omitted)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
