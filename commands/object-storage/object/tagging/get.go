package objecttagging

import (
	"context"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

func GetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object-tagging",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the tag set of an object",
		LongDesc:  "Read the key/value tag set currently attached to an object. On a versioning-enabled bucket, pass --version-id to read the tags of a specific version instead of the current one.",
		Example:   "ionosctl object-storage object tagging get --name my-bucket --key my-object",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			req := client.MustObjectStorage().ObjectStorageClient.TaggingApi.GetObjectTagging(c.Context, name, key)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			result, _, err := req.Execute()
			if err != nil {
				return err
			}

			var tags []tagInfo
			for _, t := range result.GetTagSet() {
				tags = append(tags, tagInfo{
					Key:   t.GetKey(),
					Value: t.GetValue(),
				})
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, tags, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object whose tags to read", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Read the tags of this specific object version instead of the current one (versioned buckets only)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
