package objecttagging

import (
	"context"
	"fmt"

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
		Resource:  "object-tagging",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Remove all tags from an object",
		LongDesc:  "Remove the entire tag set from an object, leaving it with no tags. There is no way to delete a single tag; to drop one tag while keeping others, use \"tagging put\" with the remaining tags. On a versioning-enabled bucket, pass --version-id to clear the tags of a specific version.",
		Example:   "ionosctl object-storage object tagging delete --name my-bucket --key my-object\nionosctl object-storage object tagging delete --name my-bucket --key my-object -f",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete tagging configuration for object %q in bucket %q", key, name), viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}

			req := client.MustObjectStorage().ObjectStorageClient.TaggingApi.DeleteObjectTagging(c.Context, name, key)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			_, _, err := req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Tagging configuration for object %q in bucket %q deleted successfully\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object whose tags to remove", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Clear the tags of this specific object version instead of the current one (versioned buckets only)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
