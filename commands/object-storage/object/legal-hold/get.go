package legalhold

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
		Resource:  "object-legal-hold",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get an object's legal hold status (ON or OFF)",
		LongDesc:  "Read whether a legal hold is currently ON or OFF for an object. On versioned buckets, pass --version-id to read a specific version's status.",
		Example:   "ionosctl object-storage object legal-hold get --name my-bucket --key my-object",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			req := client.MustObjectStorage().ObjectStorageClient.ObjectLockApi.
				GetObjectLegalHold(c.Context, name, key)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			result, _, err := req.Execute()
			if err != nil {
				return err
			}

			info := legalHoldInfo{
				Status: result.GetStatus(),
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object whose legal hold status to read", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Read the legal hold of this specific object version instead of the current one (versioned buckets only)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
