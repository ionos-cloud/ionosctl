package versioning

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

type versioningResult struct {
	Name       string `json:"Name"`
	Versioning string `json:"Versioning"`
}

func GetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "bucket-versioning",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the versioning state of a bucket",
		Example:   "ionosctl object-storage bucket versioning get --name my-bucket",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))

			s3Regional, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			apiResult, _, err := s3Regional.VersioningApi.GetBucketVersioning(c.Context, name).Execute()
			if err != nil {
				return err
			}

			status := "Disabled"
			if apiResult.HasStatus() {
				status = string(apiResult.GetStatus())
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

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
