package retention

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
		Resource:  "object-retention",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Apply a retention configuration to an object",
		LongDesc: "Place an Object Lock retention configuration on an object. " +
			"Requires the bucket to have been created with Object Lock enabled.",
		Example: "ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z\n" +
			"ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z --bypass-governance-retention",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, flagMode, flagRetainUntilDate)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			mode := viper.GetString(core.GetFlagName(c.NS, flagMode))
			retainUntilDate := viper.GetString(core.GetFlagName(c.NS, flagRetainUntilDate))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))
			bypassGovernance := viper.GetBool(core.GetFlagName(c.NS, flagBypassGovernanceRetention))

			retReq := objectstorage.NewPutObjectRetentionRequest()
			retReq.SetMode(mode)
			retReq.SetRetainUntilDate(retainUntilDate)

			req := client.MustObjectStorage().ObjectStorageClient.ObjectLockApi.
				PutObjectRetention(c.Context, name, key).
				PutObjectRetentionRequest(*retReq)
			if versionId != "" {
				req = req.VersionId(versionId)
			}
			if bypassGovernance {
				req = req.XAmzBypassGovernanceRetention(true)
			}

			_, err := req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Retention configuration for %q in bucket %q applied successfully\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagMode, "", "", "Retention mode: GOVERNANCE or COMPLIANCE", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GOVERNANCE", "COMPLIANCE"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagRetainUntilDate, "", "", "Date until which the object is retained (RFC 3339 format, e.g. 2026-01-01T00:00:00Z)", core.RequiredFlagOption())
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID of the object")
	cmd.AddBoolFlag(flagBypassGovernanceRetention, "", false, "Bypass Governance-mode restrictions")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
