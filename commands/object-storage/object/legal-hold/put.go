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
		ShortDesc: "Apply or remove a legal hold on an object",
		LongDesc: "Apply or remove a legal hold configuration on an object. " +
			"Requires the bucket to have been created with Object Lock enabled.",
		Example: "ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status ON\n" +
			"ionosctl object-storage object legal-hold put --name my-bucket --key my-object --status OFF",
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagStatus, "", "", "Legal hold status: ON or OFF", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ON", "OFF"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID of the object")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
