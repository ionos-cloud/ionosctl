package policy

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type policyStatusInfo struct {
	Bucket   string `json:"Bucket"`
	IsPublic bool   `json:"IsPublic"`
}

func StatusCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "policy",
		Verb:      "status",
		Aliases:   []string{"s"},
		ShortDesc: "Check if a bucket policy makes the bucket public",
		Example:   "ionosctl object-storage policy status --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			result, _, err := s3.PolicyApi.GetBucketPolicyStatus(c.Context, name).Execute()
			if err != nil {
				return err
			}

			info := policyStatusInfo{
				Bucket:   name,
				IsPublic: result.GetIsPublic(),
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(statusCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
