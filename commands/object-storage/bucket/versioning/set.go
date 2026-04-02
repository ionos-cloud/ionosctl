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
		Example:   "ionosctl object-storage bucket versioning set --name my-bucket --status Enabled\nionosctl object-storage bucket versioning set --name my-bucket --status Suspended",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagStatus)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			status := viper.GetString(core.GetFlagName(c.NS, flagStatus))

			s3Regional, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			versioningStatus := objectstorage.BucketVersioningStatus(status)
			req := objectstorage.NewPutBucketVersioningRequest()
			req.SetStatus(versioningStatus)

			_, err = s3Regional.VersioningApi.PutBucketVersioning(c.Context, name).
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagStatus, "", "", "Versioning status: Enabled or Suspended", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Enabled", "Suspended"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
