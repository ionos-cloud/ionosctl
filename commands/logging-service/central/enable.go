package central

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CentralEnable() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "logging-service",
		Resource:  "central",
		Verb:      "enable",
		Aliases:   []string{"e"},
		ShortDesc: "Enable CentralLogging",
		Example:   "ionosctl logging-service central enable --location de/txl",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return enable(c, true)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
