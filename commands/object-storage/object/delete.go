package object

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
)

func DeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete an object",
		Example:   "ionosctl object-storage object delete --name my-bucket --key photos/image.jpg\nionosctl object-storage object delete --name my-bucket --key photos/image.jpg --version-id abc123 -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			promptMsg := fmt.Sprintf("delete object %q from bucket %q", key, name)
			if !confirm.FAsk(c.Command.Command.InOrStdin(), promptMsg, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			s3Regional, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			req := s3Regional.ObjectsApi.DeleteObject(context.Background(), name, key)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			_, _, err = req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Object %q deleted from bucket %q\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key to delete", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagKey, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		bucket := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName))
		return completer.ObjectKeys(bucket), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID to delete a specific version")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
