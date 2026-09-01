package legalhold

import (
	"context"
	"fmt"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
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
		Resource:  "object-legal-hold",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Turn an object's legal hold ON or OFF",
		LongDesc: `Turn the legal hold on an object ON or OFF via --status.

While the hold is ON the object (version) cannot be deleted or overwritten, with no expiry date and no bypass - it stays locked until this same command sets it OFF. This is independent of any retention lock; the object remains protected while either is active.

Requires a bucket created with Object Lock enabled (it cannot be turned on for an existing bucket). On versioned buckets, pass --version-id to hold a specific version.`,
		Example: `# Place a legal hold on an object
ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status ON

# Release the legal hold
ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status OFF`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, flagStatus)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			status := viper.GetString(core.GetFlagName(c.NS, flagStatus))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))

			cfg := objectstorage.NewObjectLegalHoldConfiguration()
			cfg.SetStatus(status)

			req := client.MustObjectStorage().ObjectStorageClient.ObjectLockApi.
				PutObjectLegalHold(c.Context, name, key).
				ObjectLegalHoldConfiguration(*cfg)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			_, err := req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Legal hold %s for %q in bucket %q applied successfully\n", status, key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the Object-Lock-enabled bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object to hold or release", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagStatus, "", "", "ON to place the legal hold (locks the object indefinitely), OFF to release it", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ON", "OFF"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagVersionId, "", "", "Apply the hold to this specific object version instead of the current one (versioned buckets only)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
