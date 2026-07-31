package versioning

import (
	"context"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

const flagStatus = "status"

func SetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket-versioning",
		Verb:      "set",
		Aliases:   []string{"s"},
		ShortDesc: "Enable or suspend versioning on a bucket",
		LongDesc: `Set a bucket's versioning status to Enabled or Suspended.

Enabled starts keeping a distinct version for every overwrite/delete, protecting against accidental data loss. Suspended stops creating new versions from that point on but does NOT delete versions already stored, and does not return the bucket to a truly unversioned state. There is no way to fully disable versioning once it has been enabled.

Existing versions keep incurring storage cost; pair versioning with a lifecycle NoncurrentVersionExpiration rule if you want old versions cleaned up automatically. Versioning is also required for Object Lock.`,
		Example: `# Enable versioning
ionosctl object-storage bucket versioning set --name my-bucket --status Enabled

# Suspend versioning (existing versions are retained)
ionosctl object-storage bucket versioning set --name my-bucket --status Suspended`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagStatus)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			status := viper.GetString(core.GetFlagName(c.NS, flagStatus))

			req := objectstorage.NewPutBucketVersioningRequest()
			req.SetStatus(objectstorage.BucketVersioningStatus(status))

			_, err := client.MustObjectStorage().ObjectStorageClient.VersioningApi.PutBucketVersioning(c.Context, name).
				PutBucketVersioningRequest(*req).
				Execute()
			if err != nil {
				return err
			}

			result := versioningResult{
				Name:       name,
				Versioning: status,
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, result, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagStatus, "", "", "Target versioning status. 'Enabled' keeps a version per overwrite/delete; 'Suspended' stops new versions but retains existing ones (cannot be fully disabled once enabled)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Enabled", "Suspended"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
