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
		ShortDesc: "Set or extend the WORM retention lock on an object",
		LongDesc: `Place an Object Lock retention lock on an object, protecting it from deletion and overwrite until --retain-until-date.

--mode selects how strict the lock is: GOVERNANCE can be overridden by users with the bypass permission, COMPLIANCE can never be overridden before the date (see the "retention" group help).

The retain-until date is always in the future. You can freely EXTEND an existing lock (set a later date) in either mode. SHORTENING or removing a lock is only possible for GOVERNANCE mode and requires --bypass-governance-retention plus the bypass permission; COMPLIANCE-mode dates can never be reduced.

Requires a bucket created with Object Lock enabled - it cannot be turned on for an existing bucket. On versioned buckets, pass --version-id to lock a specific version.`,
		Example: `# Lock an object in GOVERNANCE mode until a date (RFC 3339)
ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z

# Shorten/replace an existing GOVERNANCE lock (requires bypass permission)
ionosctl object-storage object retention put --name my-bucket --key my-object --mode GOVERNANCE --retain-until-date 2026-01-01T00:00:00Z --bypass-governance-retention`,
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the Object-Lock-enabled bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object to lock", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagMode, "", "", "Retention strictness. GOVERNANCE: overridable by users with the bypass permission. COMPLIANCE: cannot be shortened or removed by anyone before the date", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GOVERNANCE", "COMPLIANCE"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagRetainUntilDate, "", "", "Future date until which the object stays locked, in RFC 3339 format (e.g. 2026-01-01T00:00:00Z). The lock lapses automatically once this passes", core.RequiredFlagOption())
	cmd.AddStringFlag(flagVersionId, "", "", "Apply the retention to this specific object version instead of the current one (versioned buckets only)")
	cmd.AddBoolFlag(flagBypassGovernanceRetention, "", false, "Required (with the bypass permission) to shorten or replace an existing GOVERNANCE-mode lock. No effect on COMPLIANCE mode")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
