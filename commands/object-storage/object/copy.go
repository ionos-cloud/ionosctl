package object

import (
	"context"
	"time"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

type copyObjectInfo struct {
	ETag         string `json:"ETag"`
	LastModified string `json:"LastModified"`
}

func CopyCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object",
		Verb:      "copy",
		Aliases:   []string{"cp"},
		ShortDesc: "Copy an object",
		LongDesc:  "Copy an object within or between buckets. The --copy-source must be in the format /source-bucket/source-key.",
		Example:   "ionosctl object-storage object copy --name dest-bucket --key dest-key --copy-source /source-bucket/source-key",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, flagCopySource)
		},
		CmdRun: func(c *core.CommandConfig) error {
			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			copySource := viper.GetString(core.GetFlagName(c.NS, flagCopySource))

			s3Regional, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			result, _, err := s3Regional.ObjectsApi.CopyObject(context.Background(), name, key).
				XAmzCopySource(copySource).
				Execute()
			if err != nil {
				return err
			}

			info := copyObjectInfo{
				ETag: result.GetETag(),
			}
			if lm := result.GetLastModified(); !lm.IsZero() {
				info.LastModified = lm.Format(time.RFC3339)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			return c.Out(table.Sprint(copyCols, info, cols))
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Destination bucket name", core.RequiredFlagOption())
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Destination object key", core.RequiredFlagOption())
	cmd.AddStringFlag(flagCopySource, "", "", "Source object in format /source-bucket/source-key", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
