package objectlock

import (
	"context"
	"strconv"

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
		Resource:  "object-lock",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the Object Lock configuration for a bucket",
		Example:   "ionosctl object-storage bucket object-lock get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			result, _, err := client.MustObjectStorage().ObjectStorageClient.ObjectLockApi.GetObjectLockConfiguration(c.Context, name).Execute()
			if err != nil {
				return err
			}

			info := configInfo{
				ObjectLockEnabled: result.GetObjectLockEnabled(),
			}

			if result.HasRule() {
				rule := result.GetRule()
				retention := rule.GetDefaultRetention()
				info.RetentionMode = retention.GetMode()
				if retention.HasDays() {
					info.RetentionDays = strconv.FormatInt(int64(retention.GetDays()), 10)
				}
				if retention.HasYears() {
					info.RetentionYears = strconv.FormatInt(int64(retention.GetYears()), 10)
				}
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
