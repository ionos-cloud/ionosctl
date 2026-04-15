package object

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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
		ShortDesc: "Download an object to a file",
		Example:   "ionosctl object-storage object get --name my-bucket --key photos/image.jpg\nionosctl object-storage object get --name my-bucket --key photos/image.jpg --destination ./local-image.jpg",
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

			s3Regional, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			req := s3Regional.ObjectsApi.GetObject(c.Context, name, key)
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key to download", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagKey, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		bucket := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName))
		return completer.ObjectKeys(bucket), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagDestination, "d", "", "Local file path for download (defaults to the basename of the key)")
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID of the object to download")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
