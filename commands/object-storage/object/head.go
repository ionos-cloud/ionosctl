package object

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

type headObjectInfo struct {
	Key           string `json:"Key"`
	ContentType   string `json:"ContentType"`
	ContentLength string `json:"ContentLength"`
	LastModified  string `json:"LastModified"`
	ETag          string `json:"ETag"`
}

func HeadCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "head",
		Aliases:   []string{"hd"},
		ShortDesc: "Get object metadata",
		Example:   "ionosctl object-storage object head --name my-bucket --key photos/image.jpg",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))

			s3Regional, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			_, apiResp, err := s3Regional.ObjectsApi.HeadObject(context.Background(), name, key).Execute()
			if err != nil {
				return err
			}

			info := headObjectInfo{
				Key:           key,
				ContentType:   apiResp.Header.Get("Content-Type"),
				ContentLength: apiResp.Header.Get("Content-Length"),
				LastModified:  apiResp.Header.Get("Last-Modified"),
				ETag:          apiResp.Header.Get("ETag"),
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(headCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagKey, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		bucket := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName))
		return completer.ObjectKeys(bucket), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
