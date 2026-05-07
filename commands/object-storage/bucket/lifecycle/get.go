package lifecycle

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
		Resource:  "lifecycle",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the lifecycle configuration for a bucket",
		Example:   "ionosctl object-storage bucket lifecycle get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			result, _, err := client.MustObjectStorage().ObjectStorageClient.LifecycleApi.GetBucketLifecycle(c.Context, name).Execute()
			if err != nil {
				return err
			}

			var rules []ruleInfo
			for _, r := range result.GetRules() {
				info := ruleInfo{
					ID:     r.GetID(),
					Prefix: r.GetPrefix(),
					Status: string(r.GetStatus()),
				}

				if exp := r.Expiration; exp != nil {
					info.ExpirationDays = int32PtrToStr(exp.Days)
					info.ExpirationDate = derefStr(exp.Date)
					info.ExpiredObjectDeleteMarker = boolPtrToStr(exp.ExpiredObjectDeleteMarker)
				}

				if nve := r.NoncurrentVersionExpiration; nve != nil {
					info.NoncurrentDays = int32PtrToStr(nve.NoncurrentDays)
				}

				if abort := r.AbortIncompleteMultipartUpload; abort != nil {
					info.AbortDays = int32PtrToStr(abort.DaysAfterInitiation)
				}

				rules = append(rules, info)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(allCols, rules, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
