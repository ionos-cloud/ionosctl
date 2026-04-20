package objectlock

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
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
		Resource:  "object-lock",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Apply an Object Lock configuration to a bucket",
		LongDesc: "Apply an Object Lock configuration to a bucket. " +
			"The bucket must have been created with --object-lock enabled. " +
			"Specify a default retention mode (GOVERNANCE or COMPLIANCE) and period (--days or --years, but not both).",
		Example: "ionosctl object-storage bucket object-lock put --name my-bucket --mode GOVERNANCE --days 30\n" +
			"ionosctl object-storage bucket object-lock put --name my-bucket --mode COMPLIANCE --years 1",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagMode); err != nil {
				return err
			}

			days := viper.GetInt32(core.GetFlagName(c.NS, flagDays))
			years := viper.GetInt32(core.GetFlagName(c.NS, flagYears))
			if days == 0 && years == 0 {
				return fmt.Errorf("at least one of --%s or --%s must be provided", flagDays, flagYears)
			}
			if days != 0 && years != 0 {
				return fmt.Errorf("--%s and --%s are mutually exclusive", flagDays, flagYears)
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			mode := viper.GetString(core.GetFlagName(c.NS, flagMode))
			days := viper.GetInt32(core.GetFlagName(c.NS, flagDays))
			years := viper.GetInt32(core.GetFlagName(c.NS, flagYears))

			retention := objectstorage.NewDefaultRetention()
			retention.SetMode(mode)
			if days > 0 {
				retention.SetDays(days)
			}
			if years > 0 {
				retention.SetYears(years)
			}

			rule := objectstorage.NewPutObjectLockConfigurationRequestRule()
			rule.SetDefaultRetention(*retention)

			req := objectstorage.NewPutObjectLockConfigurationRequest()
			req.SetObjectLockEnabled("Enabled")
			req.SetRule(*rule)

			xmlBytes, err := xml.Marshal(req)
			if err != nil {
				return fmt.Errorf("serializing object lock configuration: %w", err)
			}
			hash := md5.Sum(xmlBytes)
			contentMD5 := base64.StdEncoding.EncodeToString(hash[:])

			_, err = client.MustObjectStorage().ObjectStorageClient.ObjectLockApi.
				PutObjectLockConfiguration(c.Context, name).
				ContentMD5(contentMD5).
				PutObjectLockConfigurationRequest(*req).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Object Lock configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagMode, "", "", "Default retention mode: GOVERNANCE or COMPLIANCE", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GOVERNANCE", "COMPLIANCE"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(flagDays, "", 0, "Default retention period in days (mutually exclusive with --years)")
	cmd.AddInt32Flag(flagYears, "", 0, "Default retention period in years (mutually exclusive with --days)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
